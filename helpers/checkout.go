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
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/joho/godotenv"
	"github.com/vikas-gautam/golang_cicd/models"
)

var DestFolder = "/tmp/vikas"

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

var userdata models.UserData

// take user input and checkout code
func CodeCheckout(repoURL string, branchName string, DockerfilePath string) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	dockerRegistryUserID := os.Getenv("dockerRegistryUserID")
	if dockerRegistryUserID == "" {
		dockerRegistryUserID = "vikas93/"
	}

	dockerRepoName := os.Getenv("dockerRepoName")
	if dockerRepoName == "" {
		dockerRepoName = "go-cicd"
	}

	var dockerSrcPath = DestFolder + "/" + DockerfilePath

	DockerfileName := userdata.DockerfileName
	if DockerfileName == "" {
		DockerfileName = "Dockerfile"
	}

	Branch := branchName

	//clean workspace before cloning repo
	err = CleanWorkspace(DestFolder)

	//clone given repo
	_, errClone := git.PlainClone(DestFolder, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + Branch),
		Progress:      os.Stdout,
	})
	if errClone != nil {
		return errClone
	}
	fmt.Println("repo cloned")

	//docker client for image build and push
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = imageBuild(dockerRegistryUserID, dockerRepoName, DockerfileName, dockerSrcPath, cli)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	fmt.Println("docker image has been created")

	//push the docker image
	err = imagePush(dockerRegistryUserID, dockerRepoName, cli)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	fmt.Println("docker image has been pushed successfully")

	return nil
}

// build and create artifact
func imageBuild(dockerRegistryUserID string, dockerRepoName string, DockerfileName string, dockerSrcPath string, dockerClient *client.Client) error {
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
		Tags:       []string{dockerRegistryUserID + dockerRepoName},
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

//image push

func imagePush(dockerRegistryUserID string, dockerRepoName string, dockerClient *client.Client) error {
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

	tag := dockerRegistryUserID + dockerRepoName
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
