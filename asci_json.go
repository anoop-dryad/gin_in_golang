package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AsciiJson(ctx *gin.Context) {
	data := map[string]interface{}{
		"lang": "GO语言",
		"tag":  "<br>",
	}
	ctx.AsciiJSON(http.StatusOK, data)
}
