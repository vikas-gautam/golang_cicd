package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vikas-gautam/golang_cicd/controllers"
	"github.com/vikas-gautam/golang_cicd/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/go/homepage", controllers.Frontend)
	incomingRoutes.GET("/health", controllers.HealthCheck)

	baseApiPath := incomingRoutes.Group("/api")
	{
		//api without authentication /api/apiname
		baseApiPath.POST("/signup", controllers.Signup)
		baseApiPath.POST("/webhook/git", controllers.Gitwebhook)
		baseApiPath.POST("/webhook/docker", controllers.Dockerwebhook)
		baseApiPath.POST("/checkout", controllers.CodeCheckoutApi)

		secured := baseApiPath.Group("/secured").Use(middlewares.Auth)
		{
			//APIs with authentication /api/secured/apiname
			secured.POST("/registerApp", controllers.RegisterApp)
			secured.POST("/deployservice", controllers.DeployService)
			secured.DELETE("/deleteservice", controllers.DeleteService)
			secured.PUT("/updateservice", controllers.UpdateService)

		}
	}
}
