package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vikas-gautam/golang_cicd/routes"
)

func main() {
	ListenPort := os.Getenv("PORT")
	if ListenPort == "" {
		ListenPort = "9090"
	}
	router := gin.Default()
	routes.UserRoutes(router)
	router.Run(":" + ListenPort)
}
