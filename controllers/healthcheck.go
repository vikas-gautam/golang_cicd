package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthcheck api
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "application is ready to serve"})
}
