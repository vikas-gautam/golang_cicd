package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vikas-gautam/golang_cicd/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api/health", controllers.HealthCheck)
	incomingRoutes.POST("/api/checkout", controllers.CodeCheckoutApi)
	incomingRoutes.POST("/api/webhook/git", controllers.Gitwebhook)
	incomingRoutes.POST("/api/webhook/docker", controllers.Dockerwebhook)
	incomingRoutes.GET("/go/homepage", controllers.Frontend)
	incomingRoutes.POST("/api/registerApp", controllers.RegisterApp)
	incomingRoutes.POST("/api/deployservice", controllers.DeployService)
	incomingRoutes.DELETE("/api/deleteservice", controllers.DeleteService)
	incomingRoutes.PUT("/api/updateservice", controllers.UpdateService)

}
