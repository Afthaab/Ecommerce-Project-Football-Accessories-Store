package controllers

import (
	"strconv"
	"time"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	var cartdata models.Cart
	var productdata models.Product
	if c.Bind(&cartdata) != nil {
		c.JSON(400, gin.H{
			"Bad Request": "Could not bind the JSON data",
		})
		return
	}
	DB := config.DBconnect()
	// checking if the products is already existing in the cart
	result := DB.First(&models.Cart{}, "pid = ?", cartdata.Pid)
	//getting the data from the product table of that particular product
	DB.Raw("SELECT * FROM products WHERE productid = ?", cartdata.Pid).Scan(&productdata)
	if result.Error != nil {
		//Checking the quantity
		if cartdata.Quantity >= productdata.Stock {
			c.JSON(404, gin.H{
				"Message": "Out of Stock",
			})
			return
		}
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
		//if the product already exist then multiplying the quantity and price
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
	id, _ := strconv.Atoi(c.GetString("userid"))
	type cartdata struct {
		Productid   uint
		Productname string
		Quantity    uint
		Totalprice  uint
		Baseimage   string
		Price       string
	}
	var datas []cartdata
	DB := config.DBconnect()
	result := DB.Raw("select products.productname, products.productid, products.baseimage, carts.quantity, carts.price, carts.totalprice FROM carts INNER JOIN products ON products.productid=carts.pid WHERE cartid = ?", id).Scan(&datas)
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
func EditCart(c *gin.Context) {
	var cartdata models.Cart
	id, _ := strconv.Atoi(c.GetString("userid"))
	result := c.Bind(&cartdata)
	if result != nil {
		c.JSON(400, gin.H{
			"Error": result.Error(),
		})
		return
	}
	DB := config.DBconnect()
	result1 := DB.Exec("UPDATE carts set quantity = ? WHERE cartid = ? AND  pid = ?", cartdata.Quantity, id, cartdata.Pid)
	if result1.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Success",
	})

}

func DeleteCartItem(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Query("pid"))
	id, _ := strconv.Atoi(c.GetString("userid"))
	DB := config.DBconnect()
	result := DB.Delete(&models.Cart{}, " pid = ? AND cartid = ?", pid, id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully Deleted",
	})

}
func AddToWishlist(c *gin.Context) {
	var wishlistdata models.Wishlist
	uid, _ := strconv.Atoi(c.GetString("userid"))
	err := c.Bind(&wishlistdata)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Raw("INSERT into wishlists(pid,wishlistid) VALUES (?,?) returning id", wishlistdata.Pid, uid).Scan(&wishlistdata)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":     "Successfully added to wishlist",
		"Wishlist ID": wishlistdata.ID,
	})
}

type Wishlistitems struct {
	Id         uint
	Pid        uint
	Wishlistid uint
	CreatedAt  time.Time
}

func ViewWishlist(c *gin.Context) {
	var wishlistdata []Wishlistitems
	uid, _ := strconv.Atoi(c.GetString("userid"))
	DB := config.DBconnect()
	result := DB.Raw("SELECT * FROM wishlists WHERE wishlistid = ?", uid).Scan(&wishlistdata)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Wishlist ID": wishlistdata,
	})

}
func RemoveWishlist(c *gin.Context) {
	var wishlistdata models.Wishlist
	uid, _ := strconv.Atoi(c.GetString("userid"))
	err := c.Bind(&wishlistdata)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Exec("DELETE FROM wishlists WHERE pid = ? AND wishlistid = ?", wishlistdata.Pid, uid)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Success",
	})

}
