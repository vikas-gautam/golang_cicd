package models

type UserData struct {
	RepoURL        string `json:"repourl"  validate:"required"`
	Branch         string `json:"branch"`
	DockerfileName string `json:"DockerfileName"`
	DockerfilePath string `json:"dockerfilePath"`
}
