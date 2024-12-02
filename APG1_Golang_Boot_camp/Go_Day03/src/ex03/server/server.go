package server

import (
	"Interface/db"
	"Interface/types"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"

)


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


