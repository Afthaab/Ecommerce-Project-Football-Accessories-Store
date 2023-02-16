package controllers

import (
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

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
	DB := config.DBconnect()
	DB.Model(&addressdata).Where("uid = ?", id).Update("defaultadd", false)
	addressdata.Uid = uint(id)
	result := DB.Create(&addressdata)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	DB.Model(&addressdata).Where("addressid = ?", addressdata.Addressid).Update("defaultadd", true)
	c.JSON(200, gin.H{
		"Message":         "Address added succesfully",
		"Address ID":      addressdata.Addressid,
		"Default Address": addressdata.Defaultadd,
	})

}
func ShowAddress(c *gin.Context) {
	searchid, _ := strconv.Atoi(c.Query("addressid"))
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	type addressdata struct {
		Addressid  uint
		Name       string
		Phoneno    string
		Houseno    string
		Area       string
		Landmark   string
		City       string
		Pincode    string
		District   string
		State      string
		Country    string
		Defaultadd bool
	}
	var datas []addressdata
	DB := config.DBconnect()
	if searchid != 0 {
		result1 := DB.Raw("SELECT addressid,name, phoneno, houseno, area, landmark, city, pincode,district, state, country, defaultadd FROM addresses WHERE uid = ? AND addressid = ?", id, searchid).Scan(&datas)
		if result1.Error != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error.Error(),
			})
			return
		}
	} else {
		result := DB.Raw("SELECT addressid,name, phoneno, houseno, area, landmark, city, pincode,district, state, country, defaultadd FROM addresses WHERE uid = ?", id).Scan(&datas)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"User Addresses": datas,
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
		"Message":         "Successfully Updated the Address",
		"Address ID":      addressdata.Addressid,
		"Default Address": addressdata.Defaultadd,
	})

}
