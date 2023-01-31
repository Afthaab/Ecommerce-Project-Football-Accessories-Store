package controllers

import (
	"net/http"
	"strconv"

	"github.com/afthab/e_commerce/auth"
	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func Adminsignin(c *gin.Context) {
	type Signinadmin struct {
		Email    string
		Password string
	}
	var signindata Signinadmin
	var admindata models.Admin
	if c.Bind(&signindata) != nil {
		c.JSON(400, gin.H{
			"Message": "Could not bind the JSON Data",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.First(&admindata, "email LIKE ? AND password LIKE ?", signindata.Email, signindata.Password)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	str := strconv.Itoa(int(admindata.Adminid))
	token := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuth", token, 36000*24*30, "", "", false, true)
	c.JSON(200, gin.H{
		"Message": "Admin signin successful",
	})
}

func AdminSignout(c *gin.Context) {
	c.SetCookie("AdminAuth", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"Message": "Admin Successfully Signed Out",
	})
}

func Adminpanel(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "This is Admin Panel",
	})
}

