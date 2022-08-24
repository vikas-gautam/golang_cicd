package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/vikas-gautam/golang_cicd/helpers"
)

// healthcheck api
func Dockerwebhook(c *gin.Context) {

	jsonbyteData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonString := string(jsonbyteData)

	//fetch required fields from body
	imageTag := gojsonq.New().FromString(jsonString).Find("push_data.tag").(string)
	imagePullURL := gojsonq.New().FromString(jsonString).Find("repository.repo_name").(string)

	//create required fields to deploy application
	imageName := imagePullURL + ":" + imageTag
	containerName := "goApp"

	//printing in logs
	fmt.Println(imageTag)
	fmt.Println(imagePullURL)

	//deployment with pushed image (dockerhub/artifactory)
	if err = helpers.CD_CodeDeploy(imageName, containerName); err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deployment has been completed"})
	fmt.Println("deployment has been completed")
}
