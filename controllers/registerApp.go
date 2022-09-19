package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vikas-gautam/golang_cicd/helpers"
	"github.com/vikas-gautam/golang_cicd/models"
)

var appDataFromRequest models.RegisterAppData

func RegisterApp(c *gin.Context) {

	//putting json data into model struct
	if err := c.BindJSON(&appDataFromRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//validating received data
	validate := validator.New()
	err := validate.Struct(appDataFromRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Ensuring filename is same as app_name
	os.MkdirAll(helpers.FilePath, 0755)
	fileName := helpers.FilePath + appDataFromRequest.AppName + "." + "json"


	//IMP: if condition always run on true
	fmt.Println(*appDataFromRequest.NewOnboarding)
	if !*appDataFromRequest.NewOnboarding { //if new_onboarding: false

		//then append by searching file but also ensure services are not duplicate

		//check if file exists or not?
		if err := helpers.FileExistence(fileName); err != nil {
			fmt.Printf("File does not exist when new_onboarding: false\n")
			c.JSON(http.StatusInternalServerError, gin.H{"error": appDataFromRequest.AppName + " app_name has not been registered with us, please register it first via setting new_onboarding: true"})
		}

		//read the file
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Panicf("failed reading data from file: %s", err)
		}

		//append logic
		var existingAppData models.RegisterAppDataWrite
		json.Unmarshal(data, &existingAppData)
		fmt.Println(existingAppData)

		//check duplicacy of services and then append
		if !helpers.ContainsStruct(existingAppData.Services, appDataFromRequest.Services) {
			existingAppData.Services = append(existingAppData.Services, appDataFromRequest.Services...)
			fmt.Println(existingAppData.Services)

			//write file logic
			fileData, _ := json.MarshalIndent(existingAppData, "", " ")
			_ = ioutil.WriteFile(fileName, fileData, 0644)
			c.JSON(http.StatusOK, gin.H{"msg": "Given service has been registered under app_name " + appDataFromRequest.AppName})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": " Given service already exists under app_name " + appDataFromRequest.AppName})
		return
	}

	//if new_onboarding: true then create new file
	//check if intended file exists or not?
	if err := helpers.FileExistence(fileName); err != nil {
		fmt.Printf("File does not exist when new_onboarding: true\n")

		fileData, _ := json.MarshalIndent(appDataFromRequest, "", " ")
		_ = ioutil.WriteFile(fileName, fileData, 0644)
		c.JSON(http.StatusOK, gin.H{"Your application has been registered with us! You are all set to use CICD": appDataFromRequest})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": appDataFromRequest.AppName + " app_name name has already been regsitered with us, please choose different name"})

}
