package controllers

import (
	"net/http"
	"strconv"

	"github.com/afthab/e_commerce/auth"
	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func Adminsignin(c *gin.Context) {
	var signindata Signin
	if c.Bind(&signindata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{ //400
			"Message": "Could not bind the JSON Data",
		})
		return
	}
	var admindata models.User
	result = config.DB.First(&admindata, "email LIKE ?", signindata.Email).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{ //404
			"Error": result,
		})
		return
	}
	checkcredential := admindata.CheckPassword(signindata.Password)
	if checkcredential != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}
	if admindata.IsAdmin == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Admin not found",
		})
		return
	}
	str := strconv.Itoa(int(admindata.Userid))
	tokenstring, err := auth.TokenGeneration(str)
	token := tokenstring["access_token"]
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuth", token, 36000*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Message": "Admin signin successful",
	})
}

func AdminSignout(c *gin.Context) {
	c.SetCookie("AdminAuth", "", -1, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{
		"Message": "Admin Successfully Signed Out",
	})
}

func Adminpanel(c *gin.Context) {
	type data struct {
		Paymentmethod      string
		SuccessfulPayments string
	}
	var result error
	var paymentdata []data
	result = config.DB.Raw("SELECT paymentmethod, SUM(CASE WHEN paymentstatus='Successfull' THEN 1 ELSE 0 END) AS successful_payments FROM payments WHERE paymentmethod IN ('Wallet', 'Cash on Delivery', 'Razor Pay') GROUP BY paymentmethod").Scan(&paymentdata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Payment Method Details": paymentdata,
	})
	type data1 struct {
		DeliveredOrders  string
		PlacedOrders     string
		CancelledOrders  string
		ReturnedOrders   string
		ShippedOrders    string
		DispatchedOrders string
	}
	var paymentdata1 data1
	result = config.DB.Raw("SELECT SUM(CASE WHEN orderstatus='Delivered' THEN 1 ELSE 0 END) AS delivered_orders, SUM(CASE WHEN orderstatus='Order Placed' THEN 1 ELSE 0 END) AS placed_orders, SUM(CASE WHEN orderstatus='Order Cancelled' THEN 1 ELSE 0 END) AS cancelled_orders, SUM(CASE WHEN orderstatus='Returned' THEN 1 ELSE 0 END) AS returned_orders, SUM(CASE WHEN orderstatus='Shipped' THEN 1 ELSE 0 END) AS shipped_orders, SUM(CASE WHEN orderstatus='Order Dispatched' THEN 1 ELSE 0 END) AS dispatched_orders FROM orderditems INNER JOIN orders ON orderditems.oid=orders.orderid").Scan(&paymentdata1).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Order Details": paymentdata1,
	})
	type data3 struct {
		TotalAmount    uint
		CouponsApplied uint
	}
	var paymentdata3 data3
	result = config.DB.Raw("select sum(totalamount) as total_amount,count(couponid) as coupons_applied from payments where paymentstatus='Successfull'").Scan(&paymentdata3).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Total Amount":          paymentdata3.TotalAmount,
		"No of Coupons Applied": paymentdata3.CouponsApplied,
	})

}
