package controllers

import (
	"net/http"
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
	// DB := config.DBconnect()
	if searchid == 0 {
		result = config.DB.Raw("SELECT userid,firstname,lastname,email,phone,isblocked FROM users ORDER BY userid ASC").Scan(&userdata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
	} else {
		result = config.DB.Raw("SELECT userid,firstname,lastname,email,phone,isblocked FROM users WHERE userid = ?", searchid).Scan(&userdata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"user": userdata})

}

func UserManagement(c *gin.Context) {
	searchid := c.Query("userid")
	status := c.Query("status")
	var userdata Viewdata
	// DB := config.DBconnect()
	if status == "block" {
		result = config.DB.Raw("UPDATE users SET isblocked = true WHERE userid = ?", searchid).Scan(&userdata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Message": "Successfully blocked the User",
		})

	}
	if status == "unblock" {
		result = config.DB.Raw("UPDATE users SET isblocked = false WHERE userid = ?", searchid).Scan(&userdata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Message": "Successfully Unblocked the User",
		})
	}
}
