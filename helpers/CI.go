package helpers

import (
	"fmt"
	"log"
	"os"

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
func CI_CodeCheckout(repoURL string, branchName string, DockerfilePath string, imageVersion string) error {
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

	// clean workspace before cloning repo from helpers
	err = CleanWorkspace(DestFolder)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	//STAGE1: clone given repo
	_, errClone := git.PlainClone(DestFolder, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + Branch),
		Progress:      os.Stdout,
	})
	if errClone != nil {
		return errClone
	}
	fmt.Println("repo cloned")

	//docker client to talk with docker daemon from helpers
	cli, _ := DockerCommand_DockerClient()

	//STAGE2: build the docker image
	err = DockerCommand_ImageBuild(dockerRegistryUserID, dockerRepoName, imageVersion, DockerfileName, dockerSrcPath, cli)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	fmt.Println("docker image has been created")

	//STAGE3: push the docker image
	err = DockerCommand_ImagePush(dockerRegistryUserID, dockerRepoName, imageVersion, cli)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	fmt.Println("docker image has been pushed successfully")

	return nil
}
