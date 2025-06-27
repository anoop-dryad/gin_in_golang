package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type formA struct {
	NAME string `json:"name" xml:"name" binding:"required"`
}

type formB struct {
	AGE int `json:"age" xml:"age" binding:"required"`
}

func BindRequestBody(ctx *gin.Context) {
	objA := formA{}
	objB := formB{}

	if errA := ctx.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
		ctx.JSON(http.StatusOK, objA)
		return
	} else if errB := ctx.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
		ctx.JSON(http.StatusOK, objB)
		return
	} else if errC := ctx.ShouldBindBodyWith(&objA, binding.XML); errC == nil {
		// eg: <formA><name>Anoop</name></formA>
		ctx.JSON(http.StatusOK, objA)
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
}
