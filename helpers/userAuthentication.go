package helpers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vikas-gautam/golang_cicd/models"
)

func UserAuthentication(userName string, apiToken string) (string, string, error) {
	//read the loggedIn users file
	read_data, err := LoggedInUserdata()
	if err != nil {
		return "", "", err
	}

	var existingLoggedInDataList []models.LoggedInUserdata
	_ = json.Unmarshal(read_data, &existingLoggedInDataList)

	//before appending data check if user already exists? if yes return message
	usernameExists, matchedLoggedInData := CheckUsername(existingLoggedInDataList, userName)
	if !usernameExists {
		log.Println("Given user not found in database")
		msgUserNotExists := "given user doesn't exist,First Signup with this user"
		return msgUserNotExists, "", nil
	}

	//check if requested username and Hashpassword exists or not
	match := CheckPasswordHash(matchedLoggedInData, apiToken)
	fmt.Println(match)
	if !match {
		log.Println("Password is not matching with stored hash password")
		msgTokenNotValid := "Api_token is not valid, verify the token once"
		return msgTokenNotValid, "", nil
	}
	return "", "user authenticated successfully", nil
}
