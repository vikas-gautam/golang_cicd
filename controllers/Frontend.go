package controllers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// frontend api
func Frontend(c *gin.Context) {

	resp, err := http.Get("http://localhost:8000")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	stringData := string(body)
	log.Println(stringData)
	c.JSON(http.StatusOK, stringData)
}
