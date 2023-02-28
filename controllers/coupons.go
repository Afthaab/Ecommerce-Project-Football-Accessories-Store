package controllers

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	//checking if the entered date is valid or not
	now := time.Now()
	if coupondata.Expirationdate.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Please enter a valid date",
		})
		return
	}

	//creating the new coupon
	result = config.DB.Create(&coupondata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":         "Successfully created the coupon",
		"Coupon id":       coupondata.ID,
		"Expiration date": coupondata.Expirationdate,
	})

}
func ViewCoupons(c *gin.Context) {
	cid, _ := uuid.FromString(c.Query("couponid"))
	var coupondata []models.Coupon

	//checking and updating the expired coupon status to false
	result = config.DB.Exec("UPDATE coupons SET isactive = false WHERE expirationdate < NOW()").Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}

	//searching  a partiuclar coupon through id from the query params
	if cid != uuid.Nil {
		result = config.DB.Raw("SELECT * FROM coupons WHERE id = ?", cid).Scan(&coupondata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}

	} else { //showing all coupons
		result = config.DB.Raw("SELECT * FROM coupons").Scan(&coupondata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{
		"Coupons": coupondata,
	})

}
func RedeemCoupon(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	var coupondata models.Coupon
	var paymentdata models.Payment
	err := c.Bind(&coupondata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	//checking the coupon is used or not
	result = config.DB.First(&paymentdata, " couponid = ? AND useridno = ?", coupondata.ID, id).Error
	if result != nil {
		//getting the coupon data
		result = config.DB.Raw("SELECT * from coupons WHERE id = ?", coupondata.ID).Scan(&coupondata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		var amount uint

		//getting the totalprice from the cart
		result = config.DB.Raw("SELECT sum(totalprice) FROM carts WHERE cartid = ?", id).Scan(&amount).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}

		//checking the minimum amount to apply the coupon
		if amount < coupondata.Minamount {
			c.JSON(http.StatusNotFound, gin.H{
				"Minimun amount to redeem a coupon is ": coupondata.Minamount,
			})
			return
		}

		//checking the expiration of the coupon
		now := time.Now()
		if coupondata.Expirationdate.Before(now) {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": "Coupon has Expired",
			})
			return
		}

		//applying the discount
		final := amount - ((amount * coupondata.Discount) / 100)
		c.JSON(http.StatusOK, gin.H{
			"Message":           "Success",
			"discounted Amount": final,
			"Coupon ID":         coupondata.ID,
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Coupon alredy used by the user",
		})
	}

}
func EditCoupon(c *gin.Context) {
	cid, _ := uuid.FromString(c.Query("couponid"))
	var coupondata models.Coupon
	err := c.Bind(&coupondata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	//checking if the entered date is valid or not
	now := time.Now()
	if coupondata.Expirationdate.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Please enter a valid date",
		})
		return
	}
	//updating the coupon
	result = config.DB.Model(&coupondata).Where("id = ?", cid).Updates(models.Coupon{Couponame: coupondata.Couponame, Minamount: coupondata.Minamount, Discount: coupondata.Discount, Expirationdate: coupondata.Expirationdate, Isactive: coupondata.Isactive}).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
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

	//deleting the coupon
	result = config.DB.Delete(&models.Coupon{}, " id = ?", cid).Error
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
