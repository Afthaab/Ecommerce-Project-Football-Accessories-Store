package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	razorpay "github.com/razorpay/razorpay-go"
)

func CashOnDevlivery(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	amount, _ := strconv.Atoi(c.Query("amount"))
	cid, _ := uuid.Parse(c.Query("couponid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	paymentdata := models.Payment{
		Totalamount:   uint(amount),
		Paymentmethod: "Cash on Delivery",
		Useridno:      uint(id),
		Paymentstatus: "pending",
		Razorpayid:    "",
	}
	if cid != uuid.Nil {
		paymentdata.Couponid = cid
	}

	result = config.DB.Create(&paymentdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Paymennt Id": paymentdata.ID,
		"Message":     "Order Placed Succesfully",
	})

}
func Razorpay(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	cid, _ := uuid.Parse(c.Query("couponid"))
	amount, _ := strconv.Atoi(c.Query("amount"))
	// DB := config.DBconnect()
	var userdata models.User
	result = config.DB.Raw("SELECT * FROM users WHERE userid = ?", id).Scan(&userdata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}

	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET"))
	razpayvalue := amount * 100
	data := map[string]interface{}{
		"amount":   razpayvalue,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	value := body["id"]
	c.HTML(http.StatusOK, "app.html", gin.H{
		"userid":      userdata.Userid,
		"totalprice":  amount,
		"total":       razpayvalue,
		"paymentid":   value,
		"email":       userdata.Email,
		"phonenumber": userdata.Phone,
		"coupon":      cid,
	})

}
func RazorpaySuccess(c *gin.Context) {
	userid := c.Query("user_id")
	userID, _ := strconv.Atoi(userid)
	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	totalamount := c.Query("total")
	cid, _ := uuid.Parse(c.Query("coupon"))
	Rpay := models.RazorPay{
		UserID:          uint(userID),
		RazorPaymentId:  paymentid,
		Signature:       signature,
		RazorPayOrderID: orderid,
		AmountPaid:      totalamount,
	}
	// DB := config.DBconnect()
	result = config.DB.Create(&Rpay).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	totalprice, _ := strconv.Atoi(totalamount)
	id, _ := strconv.Atoi(userid)
	paymentdata := models.Payment{
		Totalamount:   uint(totalprice),
		Paymentmethod: "Razor Pay",
		Useridno:      uint(id),
		Paymentstatus: "successfull",
		Razorpayid:    paymentid,
		Couponid:      cid,
	}
	fmt.Println(cid)
	result = config.DB.Create(&paymentdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	pid := paymentdata.ID
	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"paymentid": pid,
	})

}
func Success(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Query("id"))
	c.HTML(http.StatusOK, "success.html", gin.H{
		"paymentid": pid,
	})

}

func WalletPayment(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	amount, _ := strconv.Atoi(c.Query("amount"))
	cid, _ := uuid.Parse(c.Query("couponid"))
	var walletdata models.Wallet
	result = config.DB.Raw("select * from wallets where walletid = ?", id).Scan(&walletdata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	if walletdata.Balance < uint(amount) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":          "Insufficeint Balance in the Wallet",
			"Wallet Balance": walletdata.Balance,
		})
		return
	}
	paymentdata := models.Payment{
		Totalamount:   uint(amount),
		Paymentmethod: "Wallet",
		Useridno:      uint(id),
		Paymentstatus: "Successfull",
		Razorpayid:    "",
	}
	if cid != uuid.Nil {
		paymentdata.Couponid = cid
	}
	result = config.DB.Create(&paymentdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	result = config.DB.Exec("update wallets set balance= balance - ? where walletid = ?", amount, id).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Paymennt Id": paymentdata.ID,
		"Message":     "Payment Succesfull",
	})
}
