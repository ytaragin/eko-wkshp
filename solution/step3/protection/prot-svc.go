package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func getping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "protection pong",
	})
}

func main() {
	log.Println("Protection Service starting up ")

	r := gin.Default()

	r.GET("/ping", getping)

	log.Println("Starting to listen...")
	r.Run() // listen and serve on 0.0.0.0:8080

}
