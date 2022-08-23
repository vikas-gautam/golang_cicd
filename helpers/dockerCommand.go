package helpers

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/joho/godotenv"
)

// initiate a client to talk to docker daemon
func DockerCommand_DockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return cli, nil
}

func DockerCommand_ImageBuild(dockerRegistryUserID string, dockerRepoName string, imageVersion string, DockerfileName string, dockerSrcPath string, dockerClient *client.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	fmt.Println(dockerSrcPath)

	tar, err := archive.TarWithOptions(dockerSrcPath, &archive.TarOptions{})
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: DockerfileName,
		Tags:       []string{dockerRegistryUserID + dockerRepoName + ":" + imageVersion},
		Remove:     true,
	}
	res, err := dockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = print(res.Body)
	if err != nil {
		return err
	}

	return nil

}

//image push

func DockerCommand_ImagePush(dockerRegistryUserID string, dockerRepoName string, imageVersion string, dockerClient *client.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	dockerUsername := os.Getenv("Username")
	if dockerUsername == "" {
		dockerUsername = "vikas93"
	}
	dockerPassword := os.Getenv("Password")

	DockerServerAddress := os.Getenv("ServerAddress")
	if DockerServerAddress == "" {
		DockerServerAddress = "https://index.docker.io/v1/"
	}

	var authConfig = types.AuthConfig{
		Username:      dockerUsername,
		Password:      dockerPassword,
		ServerAddress: DockerServerAddress,
	}

	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	tag := dockerRegistryUserID + dockerRepoName + ":" + imageVersion
	opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	rd, err := dockerClient.ImagePush(ctx, tag, opts)
	if err != nil {
		return err
	}

	defer rd.Close()

	err = print(rd)
	if err != nil {
		return err
	}

	return nil
}

// docker image pull

func DockerCommand_ImagePull(ctx context.Context, imageName string, dockerClient *client.Client) {

	out, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
}

//Container create

func DockerCommand_ContainerCreate(ctx context.Context, imageName string, containerName string, dockerClient *client.Client) (container.ContainerCreateCreatedBody, error) {

	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName}, nil, nil, nil, containerName)
	if err != nil {
		panic(err)
	}
	return resp, nil
}

//Container Start

func DockerCommand_ContainerStart(ctx context.Context, containerID string, dockerClient *client.Client) {
	if err := dockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

}

func DockerCommand_ListContainers(ctx context.Context, dockerClient *client.Client) []types.Container {
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	return containers

}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
