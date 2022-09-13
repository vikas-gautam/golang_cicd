package models

type DeployService struct {
	AppName     string `json:"app_name" validate:"required"`
	ServiceName string `json:"service_name" validate:"required"`
	UserName    string `json:"username" validate:"required"`
	ApiToken    string `json:"api_token" validate:"required"`
}
