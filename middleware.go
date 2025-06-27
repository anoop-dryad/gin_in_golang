package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		// set some value
		ctx.Set("middleware", "Logger")

		// Before request
		ctx.Next()

		// After request
		latency := time.Since(t)
		log.Println("Latency : ", latency)

		status := ctx.Writer.Status()
		log.Println("Http Status : ", status)
	}
}
