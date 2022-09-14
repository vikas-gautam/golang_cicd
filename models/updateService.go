package models

type UpdateServiceData struct {
	AppName  string        `json:"app_name" validate:"required"`
	Services []ServiceName `json:"Services" validate:"required,dive"` //using dive to ensure struct's fields are validating
	UserName string        `json:"username" validate:"required"`
	ApiToken string        `json:"api_token" validate:"required"`
}
