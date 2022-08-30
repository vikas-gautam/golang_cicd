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
	"github.com/vikas-gautam/golang_cicd/models"
)

func contains(serachIn []models.ServiceName, searchingFor string) bool {
	for _, a := range serachIn {
		if a.Name == searchingFor {
			return true
		}
	}
	return false
}

var appDataFromRequest models.RegisterAppData

func RegisterApp(c *gin.Context) {

	if err := c.BindJSON(&appDataFromRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	validate := validator.New()
	err := validate.Struct(appDataFromRequest)
	if err != nil {
		// log out this error
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(*appDataFromRequest.NewOnboarding)
	filePath := "/home/vikash/go_registerdApp/"
	os.MkdirAll(filePath, 0755)
	fileName := filePath + appDataFromRequest.AppName + "." + "json"
	fmt.Println(fileName)

	//IMP if condition always run on true
	if !*appDataFromRequest.NewOnboarding { //if new_onboarding: false
		//then append by searching file but also ensure services are not duplicate

		//check if file exists or not?
		if _, err := os.Stat(fileName); err == nil {
			fmt.Printf("File exists\n")

			//read the file
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				log.Panicf("failed reading data from file: %s", err)
			}

			//append logic
			var existingAppData models.RegisterAppData
			json.Unmarshal(data, &existingAppData)
			fmt.Println(existingAppData)

			//check duplicacy of services and then append

			if !contains(existingAppData.Services, appDataFromRequest.Services[0].Name) {
				existingAppData.Services = append(existingAppData.Services, appDataFromRequest.Services...)
				fmt.Println(existingAppData.Services)

				//write file logic
				fileData, _ := json.MarshalIndent(existingAppData, "", " ")
				_ = ioutil.WriteFile(fileName, fileData, 0644)
				c.JSON(http.StatusOK, gin.H{"msg": appDataFromRequest.Services[0].Name + " service has been registered under app_name " + appDataFromRequest.AppName})

			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": appDataFromRequest.Services[0].Name + " service already exists under app_name " + appDataFromRequest.AppName})
			}

		} else {
			fmt.Printf("File does not exist when new_onboarding: false\n")
			c.JSON(http.StatusInternalServerError, gin.H{"error": appDataFromRequest.AppName + " app_name has not been registered with us, please register it first via setting new_onboarding: true"})
		}

	} else { //if new_onboarding is true

		//then create new file
		if _, err := os.Stat(fileName); err == nil {
			fmt.Printf("File exists\n")
			c.JSON(http.StatusInternalServerError, gin.H{"error": appDataFromRequest.AppName + " app_name name has already been regsitered with us, please choose different name"})
		} else {

			fmt.Printf("File does not exists when new_onboarding: true\n")
			fileData, _ := json.MarshalIndent(appDataFromRequest, "", " ")

			_ = ioutil.WriteFile(fileName, fileData, 0644)

			c.JSON(http.StatusOK, gin.H{"Your application has been registered with us! You are all set to use CICD": appDataFromRequest})

		}

	}

}
