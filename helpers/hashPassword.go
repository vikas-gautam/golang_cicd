package helpers

import (
	b64 "encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	b64encodedPassword := b64.StdEncoding.EncodeToString([]byte(password))
	return string(bytes), b64encodedPassword, err
}

func CheckPasswordHash(b64encodedPassword, hash string) bool {
	b64decodedPassword, _ := b64.StdEncoding.DecodeString(b64encodedPassword)
	err := bcrypt.CompareHashAndPassword([]byte(hash), b64decodedPassword)
	return err == nil
}
