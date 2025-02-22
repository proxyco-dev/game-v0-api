package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRoomV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func GetRoomsV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {
	router := gin.Default();
	apiV1 := router.Group("/api/v1")

	apiV1.POST("/rooms", CreateRoomV1)
	apiV1.GET("/rooms", GetRoomsV1)

	router.Run(":8080")
}
