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
	opt := (c.Query("otp"))

	if c.Bind(&datas) != nil { //Unmarshelling the Json Data
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Could not bind the JSON Data",
		})
		return
	}
	Validationerr := validate.Struct(datas) //Validating the struct using Validator Package
	if Validationerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": Validationerr.Error(),
		})
		return
	}
	var otp string
	if opt == "email" {
		//validating the email and sending otp
		otp = initializers.Otpgeneration(datas.Email)
	}
	if opt == "phone" {
		//validating the phone number
		otp = initializers.Twilio(datas.Phone)
	}
	err := datas.HashPassword(datas.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.Create(&datas).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": result, //http.StatusNotFound because record not found
		})
		return
	}
	//setting otp in the db
	config.DB.Model(&datas).Where("email LIKE ?", datas.Email).Update("otp", otp)

	//Wallet Created
	walletdata := models.Wallet{
		Balance:  0,
		Walletid: datas.Userid,
	}
	result = config.DB.Create(&walletdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}

	//success message
	c.JSON(202, gin.H{
		"Message":  "Success", //202 success but there still one more process
		"Users ID": datas.Userid,
		"Email":    datas.Email,
		"Phone No": datas.Phone,
	})

}

func Otpvalidate(c *gin.Context) {
	var userdata models.User
	if c.Bind(&userdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Could not bind the JSON Data",
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.First(&userdata, "otp LIKE ? AND email LIKE ? AND phone LIKE ?", userdata.Otp, userdata.Email, userdata.Phone).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		config.DB.Last(&userdata).Delete(&userdata)
		c.JSON(422, gin.H{
			"Error": "Wrong OTP Register Once agian", //422 Uprocessable entity
		})
		return
	}
	result = config.DB.Exec("UPDATE users set otpverified = 'true' WHERE email = ?", userdata.Email).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "New User Successfully Registered",
		"User Id": userdata.Userid,
		"Email":   userdata.Email,
	})

}

type Signin struct {
	Email    string
	Password string
}

func Usersignin(c *gin.Context) {
	var signindata Signin
	var userdata models.User
	if c.Bind(&signindata) != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Could not bind the JSON data",
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.First(&userdata, "email LIKE ?", signindata.Email).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": result,
		})
		return
	}
	checkcredential := userdata.CheckPassword(signindata.Password)
	if checkcredential != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}
	if userdata.Otpverified == false {
		result = config.DB.Exec("DELETE FROM users WHERE email LIKE ?", signindata.Email).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "User not found",
		})
		return
	}
	if userdata.Isblocked {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "This user has been blocked by the admin",
		})
		return
	}
	str := strconv.Itoa(int(userdata.Userid))
	tokenstring, err := auth.TokenGeneration(str)
	token := tokenstring["access_token"]
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", token, 36000*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"User ID": userdata.Userid,
	})

}

func UserSignout(c *gin.Context) {
	c.SetCookie("UserAuth", "", -1, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{
		"Message": "User Successfully Signed Out",
	})
}
