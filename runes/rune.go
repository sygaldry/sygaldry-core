package runes

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// Rune contains configs for executing a Rune
type Rune struct {
	Image   string
	Tty     bool
	Env     []string
	VolumeMounts []mount.Mount
}

// NewRune creates a new Rune
func NewRune(image string, env []string, volumes []mount.Mount) (Rune, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		 return Rune{}, err
	}
	volumes = append(volumes, mount.Mount{
		Type:   mount.TypeBind,
		Source: workingDir,
		Target: "/home/sygaldry/project",
	})

	return Rune{
		Image:   image,
		Env:     env,
		Tty:     true,
		VolumeMounts: volumes,
	}, nil
}

// Run executes Rune
func (r *Rune) Run() error {
	contextBackground := context.Background()

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	_, err = dockerClient.ImagePull(contextBackground, r.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	container, err := dockerClient.ContainerCreate(
		contextBackground,
		&container.Config{
			Image: r.Image,
			Tty:   r.Tty,
			Env:   r.Env,
		},
		&container.HostConfig{
			Mounts: r.VolumeMounts,
		},
		nil,
		"",
	)

	if err := dockerClient.ContainerStart(contextBackground, container.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	_, err = dockerClient.ContainerWait(contextBackground, container.ID)
	if err != nil {
		return err
	}

	containerLogsReader, err := dockerClient.ContainerLogs(contextBackground, container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, containerLogsReader)

	return nil
}
