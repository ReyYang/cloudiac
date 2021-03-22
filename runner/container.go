package runner

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	guuid "github.com/google/uuid"
)

func (task *CommitedTask) Cancel() error {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Printf("Unable to create docker client")
		panic(err)
	}

	conatinerRemoveOpts := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}
	if err := cli.ContainerRemove(context.Background(), task.ContainerId, conatinerRemoveOpts); err != nil {
		log.Printf("Unable to remove container: %s", err)
		return err
	}
	return nil
}

func (task *CommitedTask) Status() (types.ContainerJSON, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Printf("Unable to create docker client")
		panic(err)
	}

	containerInfo, err := cli.ContainerInspect(context.Background(), task.ContainerId)
	if err != nil {
		log.Printf("Failed to inspect for container id: %s ", task.ContainerId)
		panic(err)
	}
	return containerInfo, nil
}

func (cmd *Command) Create(dirMapping string) error {
	// TODO(ZhengYue): Create client with params of host info
	cli, err := client.NewEnvClient()

	if err != nil {
		log.Printf("Unable to create docker client")
		panic(err)
	}

	id := guuid.New()

	cont, err := cli.ContainerCreate(
		cmd.ContainerInstance.Context,
		&container.Config{
			Image:        cmd.Image,
			Cmd:          cmd.Commands,
			Env:          cmd.Env,
			AttachStdout: true,
			AttachStderr: true,
		},
		&container.HostConfig{
			AutoRemove: false,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: AssetPath,
					Target: "/assets",
				},
				{
					Type:   mount.TypeBind,
					Source: dirMapping,
					Target: ContainerLogFilePath,
				},
			},
		},
		nil,
		nil,
		id.String())

	cmd.ContainerInstance.ID = cont.ID
	log.Printf("Create container ID = %s", cont.ID)

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	return err
}
