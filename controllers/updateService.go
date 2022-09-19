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

var updateDataFromRequest models.UpdateServiceData

func UpdateService(c *gin.Context) {
	//putting json data into model struct
	if err := c.BindJSON(&updateDataFromRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//validating received data
	validate := validator.New()
	err := validate.Struct(updateDataFromRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if app_name exists or not
	fileName := helpers.FilePath + updateDataFromRequest.AppName + "." + "json"

	//check if file exists or not?
	if err := helpers.FileExistence(fileName); err != nil {
		fmt.Printf("File does not exist of given app_name\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateDataFromRequest.AppName + " app_name is not registered with us, please provide valid app_name"})
	}

	//first read file to fetch required inputs to perform deletion
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}

	var existingAppData models.RegisterAppData
	json.Unmarshal(data, &existingAppData)

	serviceExists, updatedExistingAppData := helpers.CheckExistenceAndReturnStruct(existingAppData.Services, updateDataFromRequest.Services)

	if !serviceExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": deployDataFromRequest.ServiceName + " has not been registered with us"})
	}

	fmt.Println(updatedExistingAppData)

	existingAppData.Services = updatedExistingAppData

	//writing updated data back to file
	fileData, _ := json.MarshalIndent(existingAppData, "", " ")
	_ = ioutil.WriteFile(fileName, fileData, 0644)

	c.JSON(http.StatusOK, gin.H{"msg": "Given service has been updated"})

}
