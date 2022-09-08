package helpers

import (
	"github.com/vikas-gautam/golang_cicd/models"
)

func UpdateService_Contains(existingAppDatatoUpdate []models.ServiceName, updateDataFromRequest []models.ServiceName) (bool, []models.ServiceName) {
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
