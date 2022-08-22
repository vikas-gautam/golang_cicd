package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/vikas-gautam/golang_cicd/helpers"
	"github.com/vikas-gautam/golang_cicd/models"
)

// healthcheck api
func Gitwebhook(c *gin.Context) {

	jsonbyteData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonString := string(jsonbyteData)

	commitId := gojsonq.New().FromString(jsonString).Find("commits.[0].id").(string)
	cloneRepoURL := gojsonq.New().FromString(jsonString).Find("repository.clone_url").(string)
	ref := gojsonq.New().FromString(jsonString).Find("ref").(string)

	imageVersion := commitId[0:7]
	branch := strings.ReplaceAll(ref, "refs/heads/", "")

	fmt.Println(branch)
	fmt.Println(cloneRepoURL)
	fmt.Println(imageVersion)

	var userdata models.UserData
	repoUrl := fmt.Sprint(cloneRepoURL)

	userdata.RepoURL = repoUrl
	userdata.Branch = branch

	if err = helpers.CodeCheckout(userdata.RepoURL, userdata.Branch, userdata.DockerfilePath, imageVersion); err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Checkout has been completed"})
}
