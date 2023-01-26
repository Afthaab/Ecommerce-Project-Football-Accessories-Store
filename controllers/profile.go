package controllers

import (
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
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userdata Profiledata
	DB := config.DBconnect()
	result := DB.Raw("SELECT firstname,lastname,email,phone FROM users WHERE userid =?", id).Scan(&userdata)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Profile Details": userdata,
	})

}

func EditUserProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userdata models.User
	if c.Bind(&userdata) != nil {
		c.JSON(400, gin.H{
			"Error": "Unable to Bind JSON data",
		})
		return
	}
	userdata.Userid = uint(id)
	DB := config.DBconnect()
	result := DB.Model(&userdata).Updates(models.User{Firstname: userdata.Firstname, Lastname: userdata.Lastname, Phone: userdata.Phone})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Profile Updated Successfully",
	})

}

func AdminProfilepage(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("adminid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var admindata models.Admin
	DB := config.DBconnect()
	result := DB.Raw("SELECT firstname,lastname,email,phone FROM admins WHERE adminid = ?", id).Scan(&admindata)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Admin Details": admindata,
	})
}
