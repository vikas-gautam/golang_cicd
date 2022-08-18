package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// healthcheck api
func Gitwebhook(c *gin.Context) {
	payload := make(map[string]interface{})
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println(err.Error())
	}
	json, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
}

