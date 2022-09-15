package helpers

import (
	"context"
	"fmt"
	"time"
)

// to deploy code on docker container
func CD_CodeDeploy(imageName string, containerName string) error {

	//create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	//docker cli
	cli, _ := DockerClient()

	//STAGE1- pull docker image
	ImagePull(ctx, imageName, cli)

	//STAGE2- remove existing container with same name
	StopAndRemoveContainer(ctx, cli, containerName)

	//STAGE3- create container
	containerID, _ := ContainerCreate(ctx, imageName, containerName, cli)

	//STAGE4- Start container
	ContainerStart(ctx, containerID.ID, cli)

	fmt.Println(containerID)

	//STAGE5- list containers
	containers := ListContainers(ctx, cli)
	fmt.Println(containers)

	return nil
}
