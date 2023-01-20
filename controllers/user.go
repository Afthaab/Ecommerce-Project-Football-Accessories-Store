package controllers

import (
	"net/http"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/initializers"
	"github.com/afthab/e_commerce/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Usersignup(c *gin.Context) {
	var datas models.User

	if c.Bind(&datas) != nil { //Unmarshelling the Json Data
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	validationerr := validate.Struct(datas) //Validating the struct using Validator Package
	if validationerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": validationerr.Error(),
		})
		return
	}

	//validating the email and sending otp
	otp := initializers.Otpgeneration(datas.Email)

	DB := config.DBconnect()
	result := DB.Create(&datas)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	DB.Model(&datas).Where("email LIKE ?", datas.Email).Update("otp", otp)

	//success message
	c.JSON(200, gin.H{
		"message": "Go to /signup/otpvalidate",
	})

}

func Otpvalidate(c *gin.Context) {
	type Userotp struct {
		Otp string
	}
	var datas Userotp
	var userdata models.User
	if c.Bind(&datas) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.First(&userdata, "otp LIKE ?", datas.Otp)
	// result := DB.Where("otp LIKE ?", datas.Otp).Find(&userdata)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		DB.First("otp LIKE ?", datas.Otp).Delete(&userdata)
		c.JSON(http.StatusAccepted, gin.H{
			"Error":   "Wrong OTP Register Once agian",
			"Message": "Goto /signup/otpvalidate",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"Message": "Successfull Registered",
	})

}

func Usersignin(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Successful",
	})

}
