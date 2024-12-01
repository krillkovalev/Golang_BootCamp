package server

import (
	"Interface/db"
	"Interface/types"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"

)



func HandlePlaces(ctx *gin.Context) {
	
	s, present := ctx.GetQuery("page")
	if !present {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Invalid 'page' value": s})
		return
	} 
	
	page, err := strconv.Atoi(s) 
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Invalid 'page' value": s})
		return
	}
	if page < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Invalid 'page' value": s})
		return
	}

	limit := 10
	offset := (page - 1) * limit
	places, total, err := db.GetPlaces(limit, offset)
	
	if err != nil || page > total / limit {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Invalid 'page' value": s})
		return
	}

	totalpages := (total + limit - 1) / limit
	
	response := types.Response{
		Name: 		"Places",
		Total: 		total,
		Places:     places,
		Previous:   page - 1,
		Next: 		page + 1,
		Last: 		totalpages,		
	} 
	ctx.JSON(http.StatusOK, response)
}


