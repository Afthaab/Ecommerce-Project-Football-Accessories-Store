package controllers

import (
	"net/http"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

type Profiledata struct {
	Firstname string
	Lastname  string
	Email     string
	Phone     string
}

func GetUserProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userdata Profiledata
	DB := config.DBconnect()
	result := DB.Raw("SELECT firstname,lastname,email,phone FROM users WHERE userid =?", id).Scan(&userdata)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Profile Details": userdata,
	})
}

func AddAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	name := c.PostForm("name")
	phonenum := c.PostForm("phonenumber")
	houseno := c.PostForm("houseno")
	area := c.PostForm("area")
	landmark := c.PostForm("landmark")
	city := c.PostForm("city")
	pincode := c.PostForm("pincode")
	district := c.PostForm("district")
	state := c.PostForm("state")
	country := c.PostForm("country")
	address := models.Address{
		Userid:   uint(id),
		Name:     name,
		Phoneno:  phonenum,
		Houseno:  houseno,
		Area:     area,
		Landmark: landmark,
		City:     city,
		Pincode:  pincode,
		District: district,
		State:    state,
		Country:  country,
	}
	DB := config.DBconnect()
	result := DB.Create(&address)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Address added succesfully",
	})

}
func ShowAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var addressdata []models.Address
	DB := config.DBconnect()
	result := DB.Raw("SELECT * FROM addresses WHERE userid = ?", id).Scan(&addressdata)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"User Addresses": addressdata,
	})

}

// Admin Profile fuctions
func AdminProfilepage(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("adminid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var admindata models.Admin
	DB := config.DBconnect()
	result := DB.Raw("SELECT firstname,lastname,email,phone FROM admins WHERE adminid = ?", id).Scan(&admindata)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Admin Details": admindata,
	})
}
