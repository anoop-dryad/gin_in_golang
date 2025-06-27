package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(Logger()) // activating middleware

	router.GET("/", func(ctx *gin.Context) {
		val := ctx.MustGet("middleware").(string)
		ctx.JSON(http.StatusOK, gin.H{
			"message":    "pong",
			"middleware": val,
		})
	})

	router.GET("/ascii-json", AsciiJson)
	router.POST("/req-body", BindRequestBody)
	router.GET("/query-param", BindQueryParams)
	router.GET("/path-param/:name/:id", BindPathParams)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	server.ListenAndServe()
}
