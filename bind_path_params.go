package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Company struct {
	ID   string `uri:"id" binding:"required,uuid"`
	NAME string `uri:"name" binding:"required"`
}

func BindPathParams(ctx *gin.Context) {

	var company Company
	if err := ctx.ShouldBindUri(&company); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":   company.ID,
		"name": company.NAME,
	})
}
