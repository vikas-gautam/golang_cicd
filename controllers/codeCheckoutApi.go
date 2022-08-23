package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/vikas-gautam/golang_cicd/helpers"
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
func CodeCheckoutApi(c *gin.Context) {

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

	imageVersion := "latest"

	if err := c.BindJSON(&userdata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var dockerSrcPath = DestFolder + "/" + userdata.DockerfilePath

	DockerfileName := userdata.DockerfileName
	if DockerfileName == "" {
		DockerfileName = "Dockerfile"
	}

	Branch := userdata.Branch
	if Branch == "" {
		Branch = "master"
	}

	validate := validator.New()
	err = validate.Struct(userdata)
	if err != nil {
		// log out this error
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(userdata)
	c.JSON(http.StatusOK, gin.H{"Request has been successfully taken and your request was": userdata})

	//clean workspace before cloning repo
	err = helpers.CleanWorkspace(DestFolder)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//clone given repo
	_, errClone := git.PlainClone(DestFolder, false, &git.CloneOptions{
		URL:           userdata.RepoURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + Branch),
		Progress:      os.Stdout,
	})
	if errClone != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errClone.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "repo cloned"})

	//docker client for image build and push
	cli, err := helpers.DockerCommand_DockerClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = helpers.DockerCommand_ImageBuild(dockerRegistryUserID, dockerRepoName, imageVersion, DockerfileName, dockerSrcPath, cli)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "docker image has been created"})

	//push the docker image
	err = helpers.DockerCommand_ImagePush(dockerRegistryUserID, dockerRepoName, imageVersion, cli)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "docker image has been pushed successfully"})

}
