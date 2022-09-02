package models

type DeployService struct {
	AppName       string        `json:"app_name" validate:"required"`
	ServiceName   string         `json:"service_name" validate:"required"` 
}
