package helpers

import (
	"github.com/vikas-gautam/golang_cicd/models"
)

func DeployService_Contains(existingAppData []models.ServiceName, service_name string) (bool, models.ServiceName) {
	for _, existingService := range existingAppData {
		if existingService.Name == service_name {
			return true, existingService
		}
	}
	return false, models.ServiceName{}
}
