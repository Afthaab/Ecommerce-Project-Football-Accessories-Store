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
	DB := config.DBconnect()
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
	if cid != uuid.Nil {
		result := DB.Raw("SELECT * FROM coupons WHERE id = ?", cid).Scan(&coupondata).Error
		if result != nil {
			c.JSON(404, gin.H{
				"Error": result.Error(),
			})
			return
		}

	} else {
		result := DB.Raw("SELECT * FROM coupons").Scan(&coupondata).Error
		if result != nil {
			c.JSON(404, gin.H{
				"Error": result.Error(),
			})
			return
		}

	}
	c.JSON(200, gin.H{
		"Coupons": coupondata,
	})

}
func RedeemCoupon(c *gin.Context) {
	amount, _ := strconv.Atoi(c.Query("totalamount"))
	var coupondata models.Coupon
	err := c.Bind(&coupondata)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	DB := config.DBconnect()
	result := DB.First(&models.Payment{}, " couponid = ?", coupondata.ID)
	if result.Error != nil {
		result1 := DB.Raw("SELECT * from coupons WHERE id = ?", coupondata.ID).Scan(&coupondata).Error
		if result1 != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error(),
			})
			return
		}
		exp, _ := time.Parse("2006-01-02", coupondata.Expirationdate)
		currentime := time.Now().UTC()
		if currentime.After(exp) {
			DB.Exec("UPDATE coupons SET isactive = false WHERE id = ?", coupondata.ID)
			c.JSON(400, gin.H{
				"Error": "The Coupon has expired",
			})
			return
		}
		final := amount - ((amount * coupondata.Discount) / 100)
		c.JSON(200, gin.H{

			"discount Amount": final,
			"Coupon ID":       coupondata.ID,
			"exp":             exp,
		})

	} else {
		c.JSON(400, gin.H{
			"Message": "Coupon alredy used",
		})
	}

}
func EditCoupon(c *gin.Context) {
	c.Query("couponid")
	var coupondata models.Coupon
	err := c.Bind(&coupondata)
	if err != nil {
		c.JSON(200, gin.H{
			"Error": err,
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Model(&coupondata).Updates(models.Coupon{ID: coupondata.ID, Couponame: coupondata.Couponame, Minamount: coupondata.Minamount, Discount: coupondata.Discount, Expirationdate: coupondata.Expirationdate, Isactive: coupondata.Isactive})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Coupon Updated Successfully",
	})
}
