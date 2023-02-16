package controllers

import (
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/gin-gonic/gin"
)

type Viewdata struct {
	Userid    uint
	Firstname string
	Lastname  string
	Email     string
	Phone     string
	Isblocked bool
}

func Adminviewuser(c *gin.Context) {
	searchid, _ := strconv.Atoi(c.Query("userid"))
	var userdata []Viewdata
	DB := config.DBconnect()
	if searchid == 0 {
		result := DB.Raw("SELECT userid,firstname,lastname,email,phone,isblocked FROM users ORDER BY userid ASC").Scan(&userdata)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
	} else {
		result := DB.Raw("SELECT userid,firstname,lastname,email,phone,isblocked FROM users WHERE userid = ?", searchid).Scan(&userdata)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}

	}

	c.JSON(200, gin.H{"user": userdata})

}

func UserManagement(c *gin.Context) {
	searchid := c.Query("userid")
	status := c.Query("status")
	var userdata Viewdata
	DB := config.DBconnect()
	if status == "block" {
		result := DB.Raw("UPDATE users SET isblocked = true WHERE userid = ?", searchid).Scan(&userdata)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Message": "Successfully blocked the User",
		})

	}
	if status == "unblock" {
		result := DB.Raw("UPDATE users SET isblocked = false WHERE userid = ?", searchid).Scan(&userdata)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Message": "Successfully Unblocked the User",
		})
	}
}
