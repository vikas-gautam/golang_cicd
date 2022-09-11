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

var signupDataFromRequest models.SignupData
var loggedInUsersfile = "loggedInUsersfile"

func Signup(c *gin.Context) {
	//putting json data into model struct
	if err := c.BindJSON(&signupDataFromRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//validating received data
	validate := validator.New()
	err := validate.Struct(signupDataFromRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//convert text password into hash & base64encode 
	hashPassword, b64encodedPassword, _ := helpers.HashPassword(signupDataFromRequest.Password)

	fileName := FilePath + loggedInUsersfile + "." + "json"
	
	//read the file
	read_data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	var existingLoggedInDataList []models.LoggedInUserdata
	_ = json.Unmarshal(read_data, &existingLoggedInDataList)

	//Creating new data by assigning haspassword to user
	var newSignupRequest models.LoggedInUserdata
	newSignupRequest.Hashpassword = hashPassword
	newSignupRequest.Username = signupDataFromRequest.Username

	//append logic
	existingLoggedInDataList = append(existingLoggedInDataList, newSignupRequest)

	//before writing data check if haspassword matches the encoded base64 password
	match := helpers.CheckPasswordHash(b64encodedPassword, newSignupRequest.Hashpassword)
	fmt.Println(match)
	if !match {
		log.Println("Passwords are not matching so not entitled to store locally")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Inconvenience caused is deeply regreted"})
		return
	}

	//writing hashpassword data into file
	fileData, _ := json.MarshalIndent(existingLoggedInDataList, "", " ")
	_ = ioutil.WriteFile(fileName, fileData, 0644)

	c.JSON(http.StatusOK, gin.H{"msg": "User signedUp successfully, here is your api token: " + b64encodedPassword})
}
