package models

type RegisterAppData struct {
	AppName       string        `json:"app_name" validate:"required"`
	NewOnboarding *bool         `json:"new_onboarding" validate:"required"`
	Services      []ServiceName `json:"Services" validate:"required,dive"`
}

type ServiceName struct {
	Name           string `json:"name" validate:"required"`
	RepoURL        string `json:"repourl"  validate:"required"`
	Branch         string `json:"branch"`
	DockerfileName string `json:"DockerfileName"`
	DockerfilePath string `json:"dockerfilePath"`
}
