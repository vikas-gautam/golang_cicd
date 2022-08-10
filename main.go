package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ListenPort = "9090"

type UserData struct {
	RepoURL        string `json: "repourl"`
	Branch         string `json: "branch"`
	DockerfilePath string `json: dockerfilepath`
}

func main() {
	router := gin.Default()
	router.GET("/api/health", HealthCheck)
	router.POST("/api/userdata", UserInput)
	router.Run(":" + ListenPort)
}

//required apis
// /api/userdata - POST
// /api/userdata/id - POST

//healthcheck api
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "application is ready to serve"})
}

//take user input
func UserInput(c *gin.Context) {
	var userdata UserData
	if err := c.BindJSON(&userdata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(userdata)
	c.JSON(http.StatusOK, gin.H{"Request has been successfully taken and your request was": userdata})
}
