package controllers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	razorpay "github.com/razorpay/razorpay-go"
)

func CashOnDevlivery(c *gin.Context) {
	amount, _ := strconv.Atoi(c.Query("amount"))
	id, err := strconv.Atoi(c.GetString("userid"))
	cid, _ := uuid.Parse(c.Query("couponid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	method := "Cash on Delivery"
	status := "pending"
	paymentdata := models.Payment{
		Totalamount:   uint(amount),
		Paymentmethod: method,
		Useridno:      uint(id),
		Paymentstatus: status,
	}
	if cid != uuid.Nil {
		paymentdata.Couponid = cid
	}
	DB := config.DBconnect()
	result := DB.Create(&paymentdata).Error
	if result != nil {
		c.JSON(400, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(200, gin.H{
		"Paymennt Id": paymentdata.ID,
		"Message":     "Order Placed Succesfully",
	})

}
func Razorpay(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	cid, _ := uuid.Parse(c.Query("couponid"))
	DB := config.DBconnect()
	var userdata models.User
	result := DB.Raw("SELECT * FROM users WHERE userid = ?", id).Scan(&userdata)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	var amount uint
	result1 := DB.Raw("SELECT SUM(totalprice) FROM carts WHERE cartid = ?", id).Scan(&amount)
	if result1.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
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
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	value := body["id"]
	c.HTML(200, "app.html", gin.H{
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
	fmt.Println("razorpay fuckin success")
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
	DB := config.DBconnect()
	result := DB.Create(&Rpay)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	method := "Razor Pay"
	status := "pending"
	totalprice, _ := strconv.Atoi(totalamount)
	id, _ := strconv.Atoi(userid)
	paymentdata := models.Payment{
		Totalamount:   uint(totalprice),
		Paymentmethod: method,
		Useridno:      uint(id),
		Paymentstatus: status,
		Razorpayid:    paymentid,
		Couponid:      cid,
	}
	fmt.Println(cid)
	result1 := DB.Create(&paymentdata)
	if result1.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	pid := paymentdata.ID
	c.JSON(200, gin.H{

		"status":    true,
		"paymentid": pid,
	})

}
func Success(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Query("id"))
	cid := c.Query("cid")
	fmt.Printf("Fully success assholes")
	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
		"couponid":  cid,
	})

}
