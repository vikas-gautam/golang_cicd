package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-playground/validator/v10"
)

//var PersonalAccessToken = "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"
var DestFolder = "/tmp/vikas"

//var dockerRegistryUserID = ""

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type UserData struct {
	RepoURL              string `json: "repourl"  validate:"required"`
	Branch               string `json: "branch"`
	DockerfilePath       string `json: "dockerfilepath" validate:"required"`
	dockerRegistryUserID string `json: dockerRegistryUserID validate: "required"`
}

//take user input and checkout code
func CodeCheckout(c *gin.Context) {

	dockerRegistryUserID := os.Getenv("dockerRegistryUserID")
	if dockerRegistryUserID == "" {
		dockerRegistryUserID = "vikas93"
	}
	var userdata UserData
	if err := c.BindJSON(&userdata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	validate := validator.New()
	err := validate.Struct(userdata)
	if err != nil {
		// log out this error
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(userdata)
	c.JSON(http.StatusOK, gin.H{"Request has been successfully taken and your request was": userdata})

	CleanWorkspace(DestFolder)

	_, errClone := git.PlainClone(DestFolder, false, &git.CloneOptions{
		URL:           userdata.RepoURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + userdata.Branch),
		Progress:      os.Stdout,
	})
	if errClone != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errClone.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "repo cloned"})

	//image build
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = imageBuild(dockerRegistryUserID, userdata.DockerfilePath, cli)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "docker image has been created"})

}

//To remove older workspace
func CleanWorkspace(DestFolder string) {
	if err := os.RemoveAll(DestFolder); err != nil {
		fmt.Println("not able to clean ws")
	}
}

//build and create artifact
func imageBuild(dockerRegistryUserID string, DockerfilePath string, dockerClient *client.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	tar, err := archive.TarWithOptions("/tmp/vikas/attendance/", &archive.TarOptions{})
	if err != nil {
		return err
	}

	if DockerfilePath == "" {
		DockerfilePath = "Dockerfile"
	}

	opts := types.ImageBuildOptions{
		Dockerfile: DockerfilePath,
		Tags:       []string{dockerRegistryUserID + "go-cicd"},
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
