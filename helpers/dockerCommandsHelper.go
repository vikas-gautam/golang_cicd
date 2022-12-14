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
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"
)

// initiate a client to talk to docker daemon
func DockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return cli, nil
}

func ImageBuild(dockerRegistryUserID string, dockerRepoName string, imageVersion string, DockerfileName string, dockerSrcPath string, dockerClient *client.Client) error {
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

func ImagePush(dockerRegistryUserID string, dockerRepoName string, imageVersion string, dockerClient *client.Client) (string, error) {
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
		return "", err
	}

	defer rd.Close()

	err = print(rd)
	if err != nil {
		return "", err
	}

	return tag, nil
}

// docker image pull

func ImagePull(ctx context.Context, imageName string, dockerClient *client.Client) {

	out, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
}

//Container create

func ContainerCreate(ctx context.Context, imageName string, containerName string, dockerClient *client.Client) (container.ContainerCreateCreatedBody, error) {

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8000",
	}
	containerPort, _ := nat.NewPort("tcp", "80")

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName},
		&container.HostConfig{PortBindings: portBinding},
		nil, nil, containerName)
	if err != nil {
		panic(err)
	}
	return resp, nil
}

//Container Start

func ContainerStart(ctx context.Context, containerID string, dockerClient *client.Client) {
	if err := dockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

}

//list containers

func ListContainers(ctx context.Context, dockerClient *client.Client) []types.Container {

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	return containers

}

// Stop and remove a container
func StopAndRemoveContainer(ctx context.Context, dockerClient *client.Client, containername string) error {

	if err := dockerClient.ContainerStop(ctx, containername, nil); err != nil {
		log.Printf("Unable to stop container %s: %s", containername, err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := dockerClient.ContainerRemove(ctx, containername, removeOptions); err != nil {
		log.Printf("Unable to remove container: %s", err)
		return err
	}

	return nil
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
