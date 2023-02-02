package controllers

import (
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func Payment(c *gin.Context) {
	amount, _ := strconv.Atoi(c.Query("amount"))
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	method := "Cash on Delivery"
	paymentdata := models.Payment{
		Totalamount:   uint(amount),
		Paymentmethod: method,
		Uid:           uint(id),
	}
	DB := config.DBconnect()
	result := DB.Create(&paymentdata).Error
	if result != nil {
		c.JSON(400, gin.H{
			"Error": result,
		})
		return
	}
	// pid:= paymentdata.ID

	
	c.JSON(200, gin.H{
		"Message": "Order Placed Succesfully",
	})

}
