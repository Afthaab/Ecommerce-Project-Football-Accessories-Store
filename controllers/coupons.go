package controllers

import (
	"strconv"
	"time"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func AddCoupon(c *gin.Context) {
	var coupondata models.Coupon
	err := c.Bind(&coupondata)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	//checking if the entered date is valid or not
	now := time.Now()
	if coupondata.Expirationdate.Before(now) {
		c.JSON(400, gin.H{
			"Error": "Please enter a valid date",
		})
		return
	}
	DB := config.DBconnect()
	//creating the new coupon
	result := DB.Create(&coupondata).Error
	if result != nil {
		c.JSON(400, gin.H{
			"Error": result.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":         "Successfully created the coupon",
		"Coupon id":       coupondata.ID,
		"Expiration date": coupondata.Expirationdate,
	})

}
func ViewCoupons(c *gin.Context) {
	cid, _ := uuid.FromString(c.Query("couponid"))
	var coupondata []models.Coupon
	DB := config.DBconnect()
	//checking and updating the expired coupon status to false
	result := DB.Exec("UPDATE coupons SET isactive = false WHERE expirationdate < NOW()")
	if err := result.Error; err != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	//searching  a partiuclar coupon through id from the query params
	if cid != uuid.Nil {
		result1 := DB.Raw("SELECT * FROM coupons WHERE id = ?", cid).Scan(&coupondata).Error
		if result1 != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error(),
			})
			return
		}

	} else { //showing all coupons
		result1 := DB.Raw("SELECT * FROM coupons").Scan(&coupondata).Error
		if result1 != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error(),
			})
			return
		}

	}
	c.JSON(200, gin.H{
		"Coupons": coupondata,
	})

}
func RedeemCoupon(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	var coupondata models.Coupon
	var paymentdata models.Payment
	err := c.Bind(&coupondata)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	DB := config.DBconnect()
	//checking the coupon is used or not
	result := DB.First(&paymentdata, " couponid = ?", coupondata.ID)
	if result.Error != nil {
		//getting the coupon data
		result1 := DB.Raw("SELECT * from coupons WHERE id = ?", coupondata.ID).Scan(&coupondata)
		if result1.Error != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error,
			})
			return
		}
		var amount uint
		//getting the totalprice from the cart
		result := DB.Raw("SELECT sum(totalprice) FROM carts WHERE cartid = ?", id).Scan(&amount)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		//checking the minimum amount to apply the coupon
		if amount < coupondata.Minamount {
			c.JSON(404, gin.H{
				"Minimun amount to redeem a coupon is ": coupondata.Minamount,
			})
			return
		}
		//checking the expiration of the coupon
		now := time.Now()
		if coupondata.Expirationdate.Before(now) {
			c.JSON(400, gin.H{
				"Error": "Coupon has Expired",
			})
			return
		}
		//applying the discount
		final := amount - ((amount * coupondata.Discount) / 100)
		c.JSON(200, gin.H{
			"Message":           "Success",
			"discounted Amount": final,
			"Coupon ID":         coupondata.ID,
		})

	} else {
		c.JSON(400, gin.H{
			"Message": "Coupon alredy used by the user",
		})
	}

}
func EditCoupon(c *gin.Context) {
	cid, _ := uuid.FromString(c.Query("couponid"))
	var coupondata models.Coupon
	err := c.Bind(&coupondata)
	if err != nil {
		c.JSON(200, gin.H{
			"Error": err,
		})
		return
	}
	//checking if the entered date is valid or not
	now := time.Now()
	if coupondata.Expirationdate.Before(now) {
		c.JSON(400, gin.H{
			"Error": "Please enter a valid date",
		})
		return
	}
	//updating the coupon
	DB := config.DBconnect()
	result := DB.Model(&coupondata).Where("id = ?", cid).Updates(models.Coupon{Couponame: coupondata.Couponame, Minamount: coupondata.Minamount, Discount: coupondata.Discount, Expirationdate: coupondata.Expirationdate, Isactive: coupondata.Isactive})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":          "Coupon Updated Successfully",
		"Id":               cid,
		"Coupon Name":      coupondata.Couponame,
		"Minimun Amount":   coupondata.Minamount,
		"Discounted Price": coupondata.Discount,
		"Expiration Date":  coupondata.Expirationdate,
		"Active Status":    coupondata.Isactive,
	})
}
func DeleteCoupon(c *gin.Context) {
	cid, _ := uuid.FromString(c.Query("couponid"))
	DB := config.DBconnect()
	//deleting the coupon
	result := DB.Delete(&models.Coupon{}, " id = ?", cid)
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
