package controllers

import (
	"net/http"
	"strconv"

	"github.com/afthab/e_commerce/auth"
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
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON Data",
		})
		return
	}
	Validationerr := validate.Struct(datas) //Validating the struct using Validator Package
	if Validationerr != nil {
		c.JSON(400, gin.H{
			"Error": Validationerr.Error(),
		})
		return
	}

	//validating the email and sending otp
	otp := initializers.Otpgeneration(datas.Email)
	err := datas.HashPassword(datas.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err,
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Create(&datas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"message": result.Error.Error(), //404 because record not found
		})
		return
	}
	//setting otp in the db
	DB.Model(&datas).Where("email LIKE ?", datas.Email).Update("otp", otp)
	//success message
	c.JSON(202, gin.H{
		"message": "Go to /signup/otpvalidate", //202 success but there still one more process
	})

}

func Otpvalidate(c *gin.Context) {
	type Userotp struct {
		Otp   string
		Email string
	}
	var otpdata Userotp
	var userdata models.User
	if c.Bind(&otpdata) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON Data",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.First(&userdata, "otp LIKE ? AND email LIKE ?", otpdata.Otp, otpdata.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		DB.Last(&userdata).Delete(&userdata)
		c.JSON(422, gin.H{
			"Error":   "Wrong OTP Register Once agian",
			"Message": "Goto /signup/otpvalidate", //422 Uprocessable entity
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "New User Successfully Registered",
	})

}

func Usersignin(c *gin.Context) {
	type Signinuser struct {
		Email    string
		Password string
	}
	var signindata Signinuser
	var userdata models.User
	if c.Bind(&signindata) != nil {
		c.JSON(404, gin.H{
			"Message": "Could not bind the JSON data",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.First(&userdata, "email LIKE ?", signindata.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": result.Error.Error(),
		})
		return
	}
	checkcredential := userdata.CheckPassword(signindata.Password)
	if checkcredential != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}
	if userdata.Isblocked {
		c.JSON(404, gin.H{
			"Error": "This user has been blocked by the admin",
		})
		return
	}
	str := strconv.Itoa(int(userdata.Userid))
	token := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", token, 36000*24*30, "", "", false, true)
	c.JSON(200, gin.H{
		"Status":  "Signin Successful",
		"Message": "Goto /home",
	})

}
func UserSignout(c *gin.Context) {
	c.SetCookie("UserAuth", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"Message": "User Successfully Signed Out",
	})
}
