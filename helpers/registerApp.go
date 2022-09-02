package helpers

import (
	"fmt"
	"os"

	"github.com/vikas-gautam/golang_cicd/models"
)

func RegisterApp_Contains(existingAppData, appDataFromRequest []models.ServiceName) bool {
	for _, a := range appDataFromRequest {
		for _, b := range existingAppData {
			if a.Name == b.Name {
				fmt.Println(a.Name)
				return true
			}
		}
	}
	return false
}

func RegisterApp_fileExistence(fileName string) error {
	if _, err := os.Stat(fileName); err != nil {
		return err
	}
	return nil
}
