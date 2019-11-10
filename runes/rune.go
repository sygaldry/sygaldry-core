package runes

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Rune contains configs for executing a Rune
type Rune struct {
	Image   string
	Tty     bool
	Env     []string
	Volumes []string
}

// NewRune creates a new Rune
func NewRune(image string, env []string, volumes []string) (Rune, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		return Rune{}, err
	}
	volumes = append(
		volumes,
		fmt.Sprintf(
			"%s:%s",
			workingDir,
			"/home/sygaldry/project",
		),
	)

	return Rune{
		Image:   image,
		Env:     env,
		Tty:     true,
		Volumes: volumes,
	}, nil
}

// Run executes Rune
func (r *Rune) Run() error {
	contextBackground := context.Background()

	dockerRunCommand := exec.CommandContext(contextBackground, "/bin/sh", "-c", dockerRunStringBuilder(r))
	dockerRunCommandOut, _ := dockerRunCommand.StdoutPipe()
	dockerRunCommandErr, _ := dockerRunCommand.StderrPipe()
	dockerRunCommandOutScanner := bufio.NewScanner(dockerRunCommandOut)
	dockerRunCommandErrScanner := bufio.NewScanner(dockerRunCommandErr)
	go func() {
		for dockerRunCommandOutScanner.Scan() {
			fmt.Printf("%s\n", dockerRunCommandOutScanner.Text())
		}
		for dockerRunCommandErrScanner.Scan() {
			fmt.Printf("%s\n", dockerRunCommandErrScanner.Text())
		}
	}()
	runError := dockerRunCommand.Run()
	if runError != nil {
		return runError
	}

	return nil
}

func dockerRunStringBuilder(rune *Rune) string {
	var stringBuilder strings.Builder
	containerName := fmt.Sprintf("sygaldry-%d", time.Now().UnixNano())
	stringBuilder.WriteString(
		fmt.Sprintf("docker run -i --name %s ", containerName),
	)

	for _, volume := range rune.Volumes {
		stringBuilder.WriteString(
			fmt.Sprintf("-v %s ", volume),
		)
	}

	stringBuilder.WriteString(
		fmt.Sprintf("--env %s ", strings.Join(rune.Env, " --env ")),
	)

	stringBuilder.WriteString(rune.Image)

	return stringBuilder.String()
}
