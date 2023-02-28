package controllers

import (
	"net/http"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/gin-gonic/gin"
)

func CheckOut(c *gin.Context) {
	ViewCart(c)

	id, _ := strconv.Atoi(c.GetString("userid"))
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
	result = config.DB.Raw("SELECT name, phoneno, houseno, area, landmark, city, pincode,district, state, country FROM addresses WHERE defaultadd = true AND uid = ?", id).Scan(&datas).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Default Address": datas,
	})
	var totalPrice uint
	result = config.DB.Table("carts").Where("cartid = ?", id).Select("SUM(totalprice)").Scan(&totalPrice).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Total Price": totalPrice,
	})

}
