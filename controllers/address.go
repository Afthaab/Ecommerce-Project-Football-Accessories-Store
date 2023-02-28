package controllers

import (
	"net/http"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

var result error

func AddAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	var addressdata models.Address
	if c.Bind(&addressdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //400
			"Error": "Error in Binding the JSON",
		})
	}
	result = config.DB.Model(&addressdata).Where("uid = ?", id).Update("defaultadd", false).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{ //404
			"Error": result,
		})
		return
	}
	addressdata.Uid = uint(id)
	result = config.DB.Create(&addressdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //400
			"Error": result,
		})
		return
	}
	result = config.DB.Model(&addressdata).Where("addressid = ?", addressdata.Addressid).Update("defaultadd", true).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{ //400
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{ //201
		"Message":         "Address added succesfully",
		"Address ID":      addressdata.Addressid,
		"Default Address": addressdata.Defaultadd,
	})

}
func ShowAddress(c *gin.Context) {
	searchid, _ := strconv.Atoi(c.Query("addressid"))
	id, _ := strconv.Atoi(c.GetString("userid"))

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
	if searchid != 0 {
		result = config.DB.Raw("SELECT addressid,name, phoneno, houseno, area, landmark, city, pincode,district, state, country, defaultadd FROM addresses WHERE uid = ? AND addressid = ?", id, searchid).Scan(&datas).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{ //404
				"Error": result,
			})
			return
		}
	} else {
		result = config.DB.Raw("SELECT addressid,name, phoneno, houseno, area, landmark, city, pincode,district, state, country, defaultadd FROM addresses WHERE uid = ?", id).Scan(&datas).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{ //404
				"Error": result,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{ //200
		"User Addresses": datas,
	})

}

func EditAddress(c *gin.Context) {
	addressid := c.Query("addressid")
	var addressdata models.Address
	if c.Bind(&addressdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //400
			"Error": "Error in binding JSON data",
		})
		return
	}
	str, _ := strconv.Atoi(addressid)
	addressdata.Addressid = uint(str)
	result = config.DB.Model(&addressdata).Updates(models.Address{Name: addressdata.Name, Phoneno: addressdata.Phoneno, Houseno: addressdata.Houseno, Area: addressdata.Area, Landmark: addressdata.Landmark, City: addressdata.City, Pincode: addressdata.Pincode, District: addressdata.District, State: addressdata.State, Country: addressdata.Country}).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //400
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{ //200
		"Message":         "Successfully Updated the Address",
		"Address ID":      addressdata.Addressid,
		"Default Address": addressdata.Defaultadd,
	})

}
