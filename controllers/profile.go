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

func AddAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var addressdata models.Address
	if c.Bind(&addressdata) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in Binding the JSON",
		})
	}
	addressdata.Userid = uint(id)
	DB := config.DBconnect()
	result := DB.Create(&addressdata)
	if result.Error != nil {
		c.JSON(500, gin.H{
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
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var addressdata []models.Address
	DB := config.DBconnect()
	result := DB.Raw("SELECT * FROM addresses WHERE userid = ?", id).Scan(&addressdata)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"User Addresses": addressdata,
	})

}

func EditAddress(c *gin.Context) {
	addressid := c.Query("addressid")
	var addressdata models.Address
	if c.Bind(&addressdata) != nil {
		c.JSON(404, gin.H{
			"Error": "Error in binding JSON data",
		})
		return
	}
	str, err := strconv.Atoi(addressid)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": err,
		})
		return
	}
	addressdata.Addressid = uint(str)
	DB := config.DBconnect()
	result := DB.Model(&addressdata).Updates(models.Address{Name: addressdata.Name, Phoneno: addressdata.Phoneno, Houseno: addressdata.Houseno, Area: addressdata.Area, Landmark: addressdata.Landmark, City: addressdata.City, Pincode: addressdata.Pincode, District: addressdata.District, State: addressdata.State, Country: addressdata.Country})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully Updated the Address",
	})

}

// Admin Profile fuctions
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
