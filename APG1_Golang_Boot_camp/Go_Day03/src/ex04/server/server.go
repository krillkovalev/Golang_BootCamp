package server

import (
	"Interface/db"
	"Interface/types"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")



func HandlePlaces(ctx *gin.Context) {
	
	l, present := ctx.GetQuery("lat")
	if !present {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	} 
	
    lat, err := strconv.ParseFloat(l, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	lo, exist := ctx.GetQuery("lon")
	if !exist {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	
	lon, err := strconv.ParseFloat(lo, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}


	limit := 3
	places, err := db.GetPlaces(limit, lat, lon)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response := types.Response{
		Name: 		"Recommendation",
		Places:     places,
	} 
	ctx.JSON(http.StatusOK, response)
}


func createToken(username string) (string, error) { 

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
	jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func LoginHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var u types.User
	fmt.Printf("The user request value %v", u)

	tokenString, err := createToken(u.Username)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		fmt.Errorf("No such username")
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

func AuthMiddleware() gin.HandlerFunc{
	
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		ctx.Header("Content-Type", "application-json")
		if tokenString == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			fmt.Print("Missing auth header")
			return
		}
		tokenString = tokenString[len("Bearer "):]

		err := verifyToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token"})
			fmt.Print("Invalid token")
		}

		ctx.Next()
	} 
	
	

}