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

var deleteDataFromRequest models.DeployService

func DeleteService(c *gin.Context) {

	//putting json data into model struct
	if err := c.BindJSON(&deleteDataFromRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//validating received data
	validate := validator.New()
	err := validate.Struct(deleteDataFromRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Getting headers from request
	ApiToken := c.GetHeader("api_token")
	UserName := c.GetHeader("username")

	// user authentication
	validationMsg, successMsg, err := helpers.UserAuthentication(UserName, ApiToken)

	if err != nil {
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
	fileName := helpers.FilePath + deleteDataFromRequest.AppName + "." + "json"

	//check if file exists or not?
	if err := helpers.FileExistence(fileName); err != nil {
		fmt.Printf("File does not exist of given app_name\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteDataFromRequest.AppName + " app_name is not registered with us, please provide valid app_name"})
	}

	//first read file to fetch required inputs to perform deletion
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}

	var existingAppData models.RegisterAppData
	json.Unmarshal(data, &existingAppData)

	serviceExists, updatedServiceList := helpers.CheckExistenceAndDeleteService(existingAppData.Services, deleteDataFromRequest.ServiceName)
	if !serviceExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": deleteDataFromRequest.ServiceName + " has not been registered with us, provide valid service_name"})
		return
	}

	existingAppData.Services = updatedServiceList

	//printing updated data after deletion
	fmt.Println(existingAppData)

	//writing updated data back to file
	fileData, _ := json.MarshalIndent(existingAppData, "", " ")
	_ = ioutil.WriteFile(fileName, fileData, 0644)

	c.JSON(http.StatusOK, gin.H{"msg": "Given service has been removed from app"})

}
