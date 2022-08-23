package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
)

// to deploy code on docker container
func CD_CodeDeploy(imageName string, containerName string) error {

	//create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	//docker cli
	cli, _ := DockerCommand_DockerClient()

	//STAGE1- pull docker image
	DockerCommand_ImagePull(ctx, imageName, cli)

	//STAGE2- create container
	containerID, _ := DockerCommand_ContainerCreate(ctx, imageName, containerName, cli)

	//STAGE3- Start container
	DockerCommand_ContainerStart(ctx, containerID.ID, cli)

	fmt.Println(containerID)

	return nil
}

func CD_ListContainers() []types.Container {
	//create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	//docker cli
	cli, _ := DockerCommand_DockerClient()

	containers := DockerCommand_ListContainers(ctx, cli)
	return containers

}
