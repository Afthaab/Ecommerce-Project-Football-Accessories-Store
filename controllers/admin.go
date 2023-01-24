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
	str := strconv.Itoa(int(admindata.Adminid))
	token := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuth", token, 36000*24*30, "", "", false, true)
	c.JSON(http.StatusFound, gin.H{
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

type Viewdata struct {
	Userid    uint
	Firstname string
	Lastname  string
	Email     string
	Phone     string
	Isblocked bool
}

func Adminviewuser(c *gin.Context) {
	var userdata []Viewdata
	DB := config.DBconnect()
	DB.Raw("SELECT userid,firstname,lastname,email,phone,isblocked FROM users ORDER BY userid ASC").Scan(&userdata)
	c.JSON(200, gin.H{"user": userdata})

}

type Admindata struct {
	Userid string
}

func Adminsearchuser(c *gin.Context) {
	var searchdata Admindata
	var userdata Viewdata
	if c.Bind(&searchdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad Request",
		})
	}
	DB := config.DBconnect()
	result := DB.Raw("SELECT userid,firstname,lastname,email,phone,isblocked FROM users WHERE userid = ?", searchdata.Userid).Scan(&userdata)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"user": userdata,
	})

}
func Adminblockuser(c *gin.Context) {
	var userdata Viewdata
	var searchdata Admindata
	if c.Bind(&searchdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Raw("UPDATE users SET isblocked = true WHERE userid = ?", searchdata.Userid).Scan(&userdata)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"Message": "Success",
	})

}
func Adminunblockuser(c *gin.Context) {
	var userdata Viewdata
	var searchdata Admindata
	if c.Bind(&searchdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Raw("UPDATE users SET isblocked = false WHERE userid = ?", searchdata.Userid).Scan(&userdata)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"Message": "Success",
	})

}
