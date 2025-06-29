package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	router := gin.Default()
	router.Use(Logger()) // activating middleware
	router.Use(ErrorHandler())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
		v.RegisterValidation("mobile", mobileNumberValidator)
	}

	v1 := router.Group("/v1")
	v2 := router.Group("/v2") // all sync calls goes to group v2

	v1.GET("/", func(ctx *gin.Context) {
		val := ctx.MustGet("middleware").(string)
		ctx.JSON(http.StatusOK, gin.H{
			"message":    "pong",
			"middleware": val,
		})
	})

	v1.GET("/ascii-json", AsciiJson)
	v1.POST("/req-body", BindRequestBody)
	v1.GET("/query-param", BindQueryParams)
	v1.GET("/path-param/:name/:id", BindPathParams)
	v1.POST("/hotel/booking", HotelBooking)
	v1.GET("/errorhandler/middleware", func(ctx *gin.Context) {
		ctx.Error(errors.New("something went wrong"))
	})

	// async calls should use the copy of original context, should use only readonly copy.
	v2.GET("/goroutine", func(ctx *gin.Context) {
		copy_context := ctx.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			log.Println("Done! in path " + copy_context.Request.URL.Path)
		}()
		ctx.JSON(http.StatusOK, gin.H{
			"message": "DONE",
		})
	})

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")

}
