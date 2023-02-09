package controllers

import (
	"fmt"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var cartdata models.Cart
	var productdata models.Product
	if c.Bind(&cartdata) != nil {
		c.JSON(400, gin.H{
			"Bad Request": "Could not bind the JSON data",
		})
		return
	}
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	DB := config.DBconnect()
	result := DB.First(&models.Cart{}, "pid = ?", cartdata.Pid)
	DB.Raw("SELECT * FROM products WHERE productid = ?", cartdata.Pid).Scan(&productdata)
	if result.Error != nil {
		if cartdata.Quantity >= productdata.Stock {
			c.JSON(404, gin.H{
				"Message": "Out of Stock",
			})
			return
		}
		fmt.Println(productdata.Price)
		fmt.Println(cartdata.Quantity)
		totalprice := productdata.Price * cartdata.Quantity
		cartitems := models.Cart{
			Pid:        cartdata.Pid,
			Quantity:   cartdata.Quantity,
			Price:      productdata.Price,
			Totalprice: totalprice,
			Cartid:     uint(id),
		}
		result1 := DB.Create(&cartitems)
		if result1.Error != nil {
			c.JSON(400, gin.H{
				"Error": result1.Error.Error(),
			})
			return
		}
	} else {
		DB.Exec("UPDATE carts SET quantity = quantity + ? WHERE pid = ?", cartdata.Quantity, cartdata.Pid)
		DB.Raw("SELECT * FROM carts WHERE pid = ?", cartdata.Pid).Scan(&cartdata)
		totalprice := productdata.Price * cartdata.Quantity
		DB.Exec("UPDATE carts SET totalprice = ? WHERE pid = ?", totalprice, cartdata.Pid)

	}
	c.JSON(200, gin.H{
		"Message": "Added to the Cart Successfull",
	})

}

func ViewCart(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	type cartdata struct {
		Productname string
		Quantity    uint
		Totalprice  uint
		Image       string
		Price       string
	}
	var datas []cartdata
	DB := config.DBconnect()
	result := DB.Raw("select products.productname, carts.quantity, carts.price, carts.totalprice FROM carts INNER JOIN products ON products.productid=carts.pid WHERE cartid = ?", id).Scan(&datas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Cart Items": datas,
	})

}
