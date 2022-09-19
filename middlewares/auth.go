package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vikas-gautam/golang_cicd/helpers"
)

func Auth(c *gin.Context) {

	// Getting headers from request
	ApiToken := c.GetHeader("api_token")
	UserName := c.GetHeader("username")

	// user authentication
	validationMsg, successMsg, err := helpers.UserAuthentication(UserName, ApiToken)

	if err != nil {
		log.Panicf("failed reading data from loggedInUsersfile: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "failed reading data from loggedInUsersfile"})
		c.Abort()
		return
	}
	if validationMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": validationMsg})
		c.Abort()
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"msg": successMsg})

}
