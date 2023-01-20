package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/afthab/e_commerce/initializers"
	"github.com/afthab/e_commerce/models"
	"github.com/golang-jwt/jwt/v4"

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
	otp := initializers.Otpgeneration(datas.Email)
	fmt.Println(otp)

	//creating a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"otp": otp,
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenstring, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	//return response
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenstring, 36000*24*30, "", "", false, true)
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
