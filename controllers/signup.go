package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vikas-gautam/golang_cicd/helpers"
	"github.com/vikas-gautam/golang_cicd/models"
)

var signupDataFromRequest models.SignupData
var LoggedInUsersfile = "loggedInUsersfile"

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

	fileName := FilePath + LoggedInUsersfile + "." + "json"

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
	newSignupRequest.Email = signupDataFromRequest.Email

	//before appending data check if user already exists? if yes return message
	usernameExists, _ := helpers.CheckUsername(existingLoggedInDataList, newSignupRequest.Username)
	if usernameExists {
		log.Println("Given user already exists")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "given user already taken,try with other username"})
		return
	}

	//if userexists == false
	//append logic
	existingLoggedInDataList = append(existingLoggedInDataList, newSignupRequest)
	//writing hashpassword data into file
	fileData, _ := json.MarshalIndent(existingLoggedInDataList, "", " ")
	_ = ioutil.WriteFile(fileName, fileData, 0644)

	c.JSON(http.StatusOK, gin.H{"msg": "User signedUp successfully, here is your api token: " + b64encodedPassword})
}
