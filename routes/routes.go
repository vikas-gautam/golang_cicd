package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vikas-gautam/golang_cicd/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api/health", controllers.HealthCheck)
	incomingRoutes.POST("/api/checkout", controllers.CodeCheckout)
	incomingRoutes.POST("/api/webhook", controllers.Gitwebhook)
}
