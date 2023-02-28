package controllers

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"Bad Request": "Could not bind the JSON data",
		})
		return
	}

	// checking if the products is already existing in the cart
	result = config.DB.First(&models.Cart{}, "pid = ?", cartdata.Pid).Error

	//getting the data from the product table of that particular product
	config.DB.Raw("SELECT * FROM products WHERE productid = ?", cartdata.Pid).Scan(&productdata)
	if result != nil {
		//Checking the quantity
		if cartdata.Quantity >= productdata.Stock {
			c.JSON(http.StatusNotFound, gin.H{
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
		result = config.DB.Create(&cartitems).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": result,
			})
			return
		}
	} else {
		//if the product already exist then multiplying the quantity and price
		config.DB.Exec("UPDATE carts SET quantity = quantity + ? WHERE pid = ?", cartdata.Quantity, cartdata.Pid)
		config.DB.Raw("SELECT * FROM carts WHERE pid = ?", cartdata.Pid).Scan(&cartdata)
		totalprice := productdata.Price * cartdata.Quantity
		config.DB.Exec("UPDATE carts SET totalprice = ? WHERE pid = ?", totalprice, cartdata.Pid)

	}
	c.JSON(http.StatusCreated, gin.H{
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
		Price       string
	}
	var datas []cartdata

	result = config.DB.Raw("select products.productname, products.productid, carts.quantity, carts.price, carts.totalprice FROM carts INNER JOIN products ON products.productid=carts.pid WHERE cartid = ?", id).Scan(&datas).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Cart Items": datas,
	})

}
func EditCart(c *gin.Context) {
	var cartdata models.Cart
	var selectdata models.Cart
	id, _ := strconv.Atoi(c.GetString("userid"))
	result = c.Bind(&cartdata)
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}

	result = config.DB.Raw("select * from carts WHERE cartid = ? AND  pid = ?", id, cartdata.Pid).Scan(&selectdata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	sum := cartdata.Quantity * selectdata.Price
	result = config.DB.Exec("UPDATE carts set quantity = ?, totalprice = ? WHERE cartid = ? AND  pid = ?", cartdata.Quantity, sum, id, cartdata.Pid).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
	})

}

func DeleteCartItem(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Query("pid"))
	id, _ := strconv.Atoi(c.GetString("userid"))
	result = config.DB.Delete(&models.Cart{}, " pid = ? AND cartid = ?", pid, id).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully Deleted",
	})

}
func AddToWishlist(c *gin.Context) {
	var wishlistdata models.Wishlist
	uid, _ := strconv.Atoi(c.GetString("userid"))
	err := c.Bind(&wishlistdata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	result = config.DB.Raw("INSERT into wishlists(pid,wishlistid) VALUES (?,?) returning id", wishlistdata.Pid, uid).Scan(&wishlistdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
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
	result = config.DB.Raw("SELECT * FROM wishlists WHERE wishlistid = ?", uid).Scan(&wishlistdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Wishlist ID": wishlistdata,
	})

}
func RemoveWishlist(c *gin.Context) {
	var wishlistdata models.Wishlist
	uid, _ := strconv.Atoi(c.GetString("userid"))
	err := c.Bind(&wishlistdata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	result = config.DB.Exec("DELETE FROM wishlists WHERE pid = ? AND wishlistid = ?", wishlistdata.Pid, uid).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
	})

}
