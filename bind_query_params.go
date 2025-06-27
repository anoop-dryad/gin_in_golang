package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	NAME string `form:"name" binding:"required"`
	AGE  string `form:"age" binding:"required"`
}

func BindQueryParams(ctx *gin.Context) {
	var person Person

	if err := ctx.ShouldBind(&person); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, person)

}
