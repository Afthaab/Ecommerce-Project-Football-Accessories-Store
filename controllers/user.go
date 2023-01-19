package controllers

import (
	"net/http"

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

	//validating the email
	initializers.Otpgeneration(datas.Email)

	c.JSON(200, gin.H{
		"message": "Go to /signup/authentication",
	})

	// DB := config.DBconnect()
	// result := DB.Create(&datas)
	// if result.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": result.Error.Error(),
	// 	})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "User created Successfully; Go to Sign In page",
	// 	})
	// }

}

func Usersignin(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Successful",
	})

}
