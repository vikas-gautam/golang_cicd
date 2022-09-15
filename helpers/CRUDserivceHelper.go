package helpers

import (
	"fmt"
	"os"

	"github.com/vikas-gautam/golang_cicd/models"
)

// DeleteService
func CheckExistenceAndDeleteService(existingAppData []models.ServiceName, service_name string) (bool, []models.ServiceName) {
	for index, existingService := range existingAppData {
		if existingService.Name == service_name {
			return true, append(existingAppData[:index], existingAppData[index+1:]...)
		}
	}
	return false, []models.ServiceName{}
}

//RegisterApp
func ContainsStruct(existingAppData, appDataFromRequest []models.ServiceName) bool {
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

// UpdateService
func CheckExistenceAndReturnStruct(existingAppDatatoUpdate []models.ServiceName, updateDataFromRequest []models.ServiceName) (bool, []models.ServiceName) {
	for index, existingService := range existingAppDatatoUpdate {
		for _, updatedServiceData := range updateDataFromRequest {
			if existingService.Name == updatedServiceData.Name {
				//Very IMP
				existingAppDatatoUpdate[index] = updatedServiceData
				return true, existingAppDatatoUpdate
			}
		}
	}
	return false, []models.ServiceName{}
}



// common fileExistence func
func FileExistence(fileName string) error {
	if _, err := os.Stat(fileName); err != nil {
		return err
	}
	return nil
}