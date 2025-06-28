package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,bookabledate" time_format:"2006-01-02"`
	MobNum   int64     `form:"mob" binding:"required,mobile"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

var mobileNumberValidator validator.Func = func(fl validator.FieldLevel) bool {
	mobNum, ok := fl.Field().Interface().(int64)
	if ok {
		if len(strconv.FormatInt(mobNum, 10)) != 10 {
			return false
		}
	}
	return true
}

func HotelBooking(ctx *gin.Context) {
	var booking Booking

	if err := ctx.ShouldBindWith(&booking, binding.Query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"checkin":  booking.CheckIn,
		"checkout": booking.CheckOut,
		"mobile":   booking.MobNum,
	})

}
