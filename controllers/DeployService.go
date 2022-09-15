package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vikas-gautam/golang_cicd/helpers"
	"github.com/vikas-gautam/golang_cicd/models"
)

var deployDataFromRequest models.DeployService

func DeployService(c *gin.Context) {

	//putting json data into model struct
	if err := c.BindJSON(&deployDataFromRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//validating received data
	validate := validator.New()
	err := validate.Struct(deployDataFromRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// user authentication
	validationMsg, successMsg, errReadingData := helpers.UserAuthentication(deployDataFromRequest.UserName, deployDataFromRequest.ApiToken)

	if errReadingData != nil {
		log.Panicf("failed reading data from loggedInUsersfile: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "failed reading data from loggedInUsersfile"})
		return
	}
	if validationMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": validationMsg})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"msg": successMsg})

	//check if app_name exists or not
	fileName := helpers.FilePath + deployDataFromRequest.AppName + "." + "json"

	//check if file exists or not?
	if err := helpers.FileExistence(fileName); err != nil {
		fmt.Printf("File does not exist of given app_name\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": deployDataFromRequest.AppName + " app_name has not been registered with us, please register it first via calling register app api"})
	}

	//call internal api to perform CICD

	//first read file to fetch required inputs for CI
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}

	var existingAppData models.RegisterAppData
	json.Unmarshal(data, &existingAppData)

	serviceExists, matchedServicedata := helpers.DeployService_Contains(existingAppData.Services, deployDataFromRequest.ServiceName)

	if !serviceExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": deployDataFromRequest.ServiceName + " has not been registered with us"})
	}
	//service_name exists so doing CICD

	//CI Stage
	imageVersion := deployDataFromRequest.AppName + "-" + deployDataFromRequest.ServiceName + "-" + "latest"
	pushedImgTag, err := helpers.CI_CodeCheckout(matchedServicedata.RepoURL, matchedServicedata.Branch, matchedServicedata.DockerfilePath, imageVersion)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": "CI has been completed"})

	//CD Stage

	//fetch required fields from body
	//Challenge is to fetch same variable those were being used while pushing image
	//imagePullURL := dockerRegistryUserID + "/" + dockerRepoName

	//create required fields to deploy application
	imageName := pushedImgTag
	containerName := "goApp" //app_name-service_name

	//deployment with pushed image (dockerhub/artifactory)
	if err = helpers.CD_CodeDeploy(imageName, containerName); err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deployment has been completed"})
	fmt.Println("deployment has been completed")

}
