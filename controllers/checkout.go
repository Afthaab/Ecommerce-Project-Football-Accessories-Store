package controllers

import (
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/gin-gonic/gin"
)

func CheckOut(c *gin.Context) {
	ViewCart(c)

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	type addressdata struct {
		Name     string
		Phoneno  string
		Houseno  string
		Area     string
		Landmark string
		City     string
		Pincode  string
		District string
		State    string
		Country  string
	}
	var datas []addressdata
	DB := config.DBconnect()
	result := DB.Raw("SELECT name, phoneno, houseno, area, landmark, city, pincode,district, state, country FROM addresses WHERE defaultadd = true AND uid = ?", id).Scan(&datas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(200, gin.H{
		"Default Address": datas,
	})
	var totalPrice float64
	result1 := DB.Table("carts").Where("cartid = ?", id).Select("SUM(totalprice)").Scan(&totalPrice).Error
	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	c.JSON(200, gin.H{
		"Total Price": totalPrice,
	})

}
