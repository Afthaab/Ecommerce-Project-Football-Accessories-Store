package controllers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Bad Request",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.First(&admindata, "email LIKE ? AND password LIKE ?", signindata.Email, signindata.Password)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"Message": "Admin signin successful",
	})
}
func Adminviewuser(c *gin.Context) {
	type viewdata struct {
		Userid    uint
		Firstname string
		Lastname  string
		Email     string
		Phone     string
		Isblocked bool
	}
	var userdata []viewdata
	DB := config.DBconnect()
	DB.Raw("SELECT userid,firstname,lastname,email,phone,isblocked FROM users ORDER BY userid ASC").Scan(&userdata)
	c.JSON(200, gin.H{"user": userdata})

}

func Adminpanel(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "This is Admin Panel",
	})
}
