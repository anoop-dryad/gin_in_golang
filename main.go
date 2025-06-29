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
	router.GET("/errorhandler/middleware", func(ctx *gin.Context) {
		ctx.Error(errors.New("something went wrong"))
	})
	// async calls should use the copy of original context, should use only readonly copy.
	router.GET("/goroutine", func(ctx *gin.Context) {
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
