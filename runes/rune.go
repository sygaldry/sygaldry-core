package runes

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Rune contains configs for executing a Rune
type Rune struct {
	Image string
	Tty   bool
	Env   []string
}

// NewRune creates a new Rune
func NewRune(image string, env []string) Rune {
	return Rune{
		Image: image,
		Env:   env,
		Tty:   true,
	}
}

// Run executes Rune
func (r *Rune) Run() error {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	_, err = cli.ImagePull(ctx, r.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: r.Image,
		Tty:   r.Tty,
		Env:   r.Env,
	}, nil, nil, "")

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	_, err = cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, out)

	return nil
}
