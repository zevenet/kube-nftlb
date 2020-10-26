package dsr

import (
	"context"
	"fmt"

	dockerTypes "github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
)

var (
	dockerCli *dockerClient.Client
)

func dockerCmdRun(action string, virtualAddr string, container string) error {
	// Link the configuration to the container where we want to apply said configuration
	// (once linked, it is only necessary to launch the "Start" command)
	exec, err := dockerCli.ContainerExecCreate(context.TODO(), container, dockerTypes.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          []string{"/bin/sh", "-c", fmt.Sprintf("ip ad %s %s/32 dev lo", action, virtualAddr)},
		Tty:          true,
		Detach:       false,
		Privileged:   true,
		User:         "root",
		WorkingDir:   "/",
	})
	if err != nil {
		return err
	}

	return dockerCli.ContainerExecStart(context.TODO(), exec.ID, dockerTypes.ExecStartCheck{
		Detach: false,
		Tty:    true,
	})
}

func init() {
	var err error
	dockerCli, err = dockerClient.NewClientWithOpts(dockerClient.FromEnv)
	if err != nil {
		panic(err)
	}
}
