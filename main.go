package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	router := gin.Default()
	router.Use(Logger()) // activating middleware

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
		v.RegisterValidation("mobile", mobileNumberValidator)
	}

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
	router.POST("/hotel/booking", HotelBooking)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	server.ListenAndServe()
}
