package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vikas-gautam/golang_cicd/models"
)

var appData models.RegisterAppData

func RegisterApp(c *gin.Context) {

	if err := c.BindJSON(&appData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	validate := validator.New()
	err := validate.Struct(appData)
	if err != nil {
		// log out this error
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(*appData.NewOnboarding)

	if !*appData.NewOnboarding { //if false
		//then append by searching file.

		//get the file name
		Name := appData.AppName + "." + "json"
		libRegEx, e := regexp.Compile(Name)
		if e != nil {
			log.Fatal(e)
		}

		//find the file name
		e = filepath.Walk("/home/vikash/goProjects/golang_cicd/", func(path string, info os.FileInfo, err error) error {
			if err == nil && libRegEx.MatchString(info.Name()) {
				println(info.Name())
			}
			return nil
		})
		if e != nil {
			log.Fatal(e)
		}

		//read the file
		data, err := ioutil.ReadFile(Name)
		if err != nil {
			log.Panicf("failed reading data from file: %s", err)
		}
		var newAppData models.RegisterAppData
		json.Unmarshal(data, &newAppData)
		fmt.Println(newAppData)
		newAppData.Services = append(newAppData.Services, appData.Services...)
		fmt.Println(newAppData.Services)
		fmt.Println(newAppData)

		fileData, _ := json.MarshalIndent(newAppData, "", " ")

		_ = ioutil.WriteFile(Name, fileData, 0644)

	} else {
		//then create new file

		fileName := appData.AppName + "." + "json"

		fileData, _ := json.MarshalIndent(appData, "", " ")

		_ = ioutil.WriteFile(fileName, fileData, 0644)

		c.JSON(http.StatusOK, gin.H{"Your application has been registered with us! You are all set to use CICD": appData})

	}

	

}
