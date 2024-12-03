package main

import (
	"Interface/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	authorized := router.Group("/api")
	authorized.Use(server.AuthMiddleware())
	authorized.GET("/recommend", server.HandlePlaces, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"response": "Auth is working!"})
	})
	
	
	router.GET("/api/get_token", server.LoginHandler)
	router.Run("127.0.0.1:8888")

}
