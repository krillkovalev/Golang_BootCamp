package main

import (
	"Interface/server"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/places", server.HandlePlaces)
	router.Run("127.0.0.1:8888")

}
