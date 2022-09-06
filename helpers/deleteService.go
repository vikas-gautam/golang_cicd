package helpers

import (
	"github.com/vikas-gautam/golang_cicd/models"
)

func DeleteService_Remove(existingAppData []models.ServiceName, service_name string) (bool, []models.ServiceName) {
	for index, existingService := range existingAppData {
		if existingService.Name == service_name {
			return true, append(existingAppData[:index], existingAppData[index+1:]...)
		}
	}
	return false, []models.ServiceName{}
}
