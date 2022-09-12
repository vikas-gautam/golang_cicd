package helpers

import (
	b64 "encoding/base64"

	"github.com/vikas-gautam/golang_cicd/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	b64encodedPassword := b64.StdEncoding.EncodeToString([]byte(password))
	return string(bytes), b64encodedPassword, err
}

func CheckPasswordHash(matchedLoggedInData models.LoggedInUserdata, b64encodedPassword string) bool {
	b64decodedPassword, _ := b64.StdEncoding.DecodeString(b64encodedPassword)
	err := bcrypt.CompareHashAndPassword([]byte(matchedLoggedInData.Hashpassword), b64decodedPassword)
	return err == nil
}

func CheckUsername(existingLoggedInDataList []models.LoggedInUserdata, findUser string) (bool, models.LoggedInUserdata) {
	for _, existingUsersData := range existingLoggedInDataList {
		if existingUsersData.Username == findUser {
			return true, existingUsersData
		}
	}
	return false, models.LoggedInUserdata{}
}
