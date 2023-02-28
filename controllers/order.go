package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func PlaceOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	pid, _ := strconv.Atoi(c.Query("paymentid"))
	aid, _ := strconv.Atoi(c.Query("addressid"))

	//getting the payment details through the payment id recieved in the query params
	var paymentdata models.Payment
	result = config.DB.Raw("SELECT * FROM payments WHERE id = ?", pid).Scan(&paymentdata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}

	//auto inserting details from the payment details to the orders
	orderdata := models.Orders{
		Useridno:    uint(id),
		Totalamount: paymentdata.Totalamount,
		Paymentid:   uint(pid),
		Orderstatus: "Order Placed",
	}

	//checking if the user has changed the address if yes then recieved through the query params
	if aid != 0 {
		orderdata.Addid = uint(aid)
	} else {
		result = config.DB.Raw("select addressid from addresses where defaultadd=true and uid = ?", id).Scan(&orderdata.Addid).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
	}
	//creating the order
	result = config.DB.Create(&orderdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	//shifting all the products in the carts to the orderditems table to save the products from the cart
	var itemdata models.Orderditems
	result = config.DB.Raw("insert into orderditems(pid, oid, quantity, price, totalprice, uid) select carts.pid, orders.orderid, carts.quantity, carts.price, carts.totalprice, orders.useridno from orders inner join carts on orders.useridno=carts.cartid where orderid = ?", orderdata.Orderid).Scan(&itemdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	//deleting the products in the cart
	result = config.DB.Exec("DELETE FROM carts WHERE cartid = ?", id).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":  "Successfullt placed order",
		"Order ID": orderdata.Orderid,
	})

}

type Orderdetails struct {
	Productname string
	Id          uint
	Pid         uint
	Oid         uuid.UUID
	Quantity    uint
	Price       uint
	Totalprice  uint
	Orderstatus string
	Useridno    uint
	Date        time.Time
}

func ViewOrder(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	prodid, _ := strconv.Atoi(c.Query("productid"))
	uid, _ := strconv.Atoi(c.GetString("userid"))
	var data []Orderdetails

	//to check if the oid or pid is passed through the query params
	if oid != uuid.Nil || prodid != 0 {

		//showing the shipment address
		var addressesdata models.Address
		result = config.DB.Raw("select addresses.* from orders inner join addresses on addresses.addressid=orders.addid where orders.orderid = ?", oid).Scan(&addressesdata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Address ID": addressesdata.Addressid,
			"Name":       addressesdata.Name,
			"House No":   addressesdata.Houseno,
			"Area":       addressesdata.Area,
			"Landmark":   addressesdata.Landmark,
			"City":       addressesdata.City,
			"Pincode":    addressesdata.Pincode,
			"District":   addressesdata.District,
			"State":      addressesdata.State,
			"Country":    addressesdata.Country,
		})

		//showing the products based on the query params
		query := "select products.productname, orderditems.*,orders.orderstatus, orders.useridno, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid where"
		if prodid == 0 {
			query = fmt.Sprintf("%s oid = '%s'", query, oid)
		} else {
			query = fmt.Sprintf("%s oid = '%s' and pid = %d", query, oid, prodid)
		}
		result = config.DB.Raw(query).Scan(&data).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Order Details": data,
		})

		//showing the payment details for the particular order
		var paymentdata models.Payment
		result = config.DB.Raw("select payments.* from orders INNER JOIN payments on orders.paymentid=payments.id where orderid = ?", oid).Scan(&paymentdata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Payment Id":     paymentdata.ID,
			"Total Amount":   paymentdata.Totalamount,
			"Payment Method": paymentdata.Paymentmethod,
			"Coupon ID":      paymentdata.Couponid,
			"Razorpay ID":    paymentdata.Razorpayid,
			"Created At":     paymentdata.CreatedAt,
			"Updated At":     paymentdata.UpdatedAt,
			"Payment Status": paymentdata.Paymentstatus,
		})

	} else {

		// if no query params then just the list of all orderd products
		result = config.DB.Raw("select products.productname, orderditems.*,orders.orderstatus, orders.useridno, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid where uid = ?", uid).Scan(&data).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Your Orders": data,
		})
	}

}

func CancelOrder(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))

	//changing the status in the orders table
	result = config.DB.Exec("UPDATE orders SET orderstatus = 'Order Cancelled' WHERE orderid = ?", oid).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}

	//to get the payment id and method to update the status in the payment table
	var paymentdata models.Payment
	result = config.DB.Raw("select payments.paymentmethod,id,payments.useridno from orders inner join payments on payments.id=orders.paymentid where orderid = ?", oid).Scan(&paymentdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	//if the payment in done in razor pay inserting the balance to the wallet and updating the status in the payment table
	if paymentdata.Paymentmethod == "Razor Pay" {
		var sum uint
		result = config.DB.Raw("select sum(totalprice) from orderditems where oid = ?", oid).Scan(&sum).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		result = config.DB.Exec("update wallets set balance = balance + ? where walletid = ?", sum, paymentdata.Useridno).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": result,
			})
			return
		}

		result = config.DB.Exec("UPDATE payments set paymentstatus = 'refunded' WHERE id = ?", paymentdata.ID).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": result,
			})
			return
		}
	} else {
		//if the payment is done in cod just updating the payment status
		result = config.DB.Exec("UPDATE payments set paymentstatus = 'Failed' WHERE id = ?", paymentdata.ID).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": result,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":   "Success",
		"Orderd ID": oid,
	})
}

func ReturnOrder(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	uid, _ := strconv.Atoi(c.GetString("userid"))
	result = config.DB.Exec("UPDATE orders set orderstatus = 'Return requested' WHERE orderid = ? AND useridno = ?", oid, uid).Error
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

func ShowWallet(c *gin.Context) {
	var walletdata models.Wallet
	uid, _ := strconv.Atoi(c.GetString("userid"))

	//Getting the wallet data
	result = config.DB.Raw("select * from wallets where walletid = ?", uid).Scan(&walletdata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Wallet ID":      walletdata.ID,
		"Wallet Balance": walletdata.Balance,
		"User ID":        walletdata.Walletid,
		"Created At":     walletdata.CreatedAt,
		"Updated At":     walletdata.UpdatedAt,
	})
}

// ------------------------ Admin Route Functions ----------------------
func AdminOrderView(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	pid, _ := strconv.Atoi(c.Query("pid"))

	var data []Orderdetails
	// to search a particular order
	if oid != uuid.Nil || pid != 0 {
		var addressesdata models.Address
		// getting the address from orders
		result = config.DB.Raw("select addresses.* from orders inner join addresses on addresses.addressid=orders.addid where orders.orderid = ?", oid).Scan(&addressesdata).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Address ID": addressesdata.Addressid,
			"Name":       addressesdata.Name,
			"House No":   addressesdata.Houseno,
			"Area":       addressesdata.Area,
			"Landmark":   addressesdata.Landmark,
			"City":       addressesdata.City,
			"Pincode":    addressesdata.Pincode,
			"District":   addressesdata.District,
			"State":      addressesdata.State,
			"Country":    addressesdata.Country,
		})

		//showing the products based on the query params
		query := "select products.productname, orderditems.*,orders.orderstatus, orders.useridno, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid where"
		if pid == 0 {
			query = fmt.Sprintf("%s oid = '%s'", query, oid)
		} else {
			query = fmt.Sprintf("%s oid = '%s' and pid = %d", query, oid, pid)
		}
		result = config.DB.Raw(query).Scan(&data).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		//showing the payment details for the particular order
		var paymentdata models.Payment
		result = config.DB.Raw("select payments.* from orders INNER JOIN payments on orders.paymentid=payments.id where orderid = ?", oid).Scan(&paymentdata).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Payment Id":     paymentdata.ID,
			"Total Amount":   paymentdata.Totalamount,
			"Payment Method": paymentdata.Paymentmethod,
			"Coupon ID":      paymentdata.Couponid,
			"Razorpay ID":    paymentdata.Razorpayid,
			"Created At":     paymentdata.CreatedAt,
			"Updated At":     paymentdata.UpdatedAt,
			"Paymemt Status": paymentdata.Paymentstatus,
		})
	} else {
		//showing all orders
		result = config.DB.Raw("select products.productname, orderditems.*,orders.orderstatus,orders.useridno, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid;").Scan(&data).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"Orders": data,
	})

}
func ReturnAccepted(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	var orderdata models.Orders
	result = config.DB.Raw("UPDATE orders SET orderstatus = 'Returned' WHERE orderid = ? RETURNING totalamount,useridno, paymentid", oid).Scan(&orderdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	result = config.DB.Exec("update wallets set balance = balance + ? where walletid = ?", orderdata.Totalamount, orderdata.Useridno).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	result = config.DB.Exec("UPDATE payments SET paymentstatus = 'Refunded' WHERE id = ?", orderdata.Paymentid).Error
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
func StatusUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("status"))
	var status string
	if id == 2 {
		status = "Shipped"
	} else if id == 3 {
		status = "Order Dispatched"
	} else if id == 4 {
		status = "Delivered"
	}
	oid, _ := uuid.FromString(c.Query("orderid"))
	var pid uint
	result = config.DB.Raw("UPDATE orders SET orderstatus = ? WHERE orderid = ? RETURNING paymentid", status, oid).Scan(&pid).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}

	if status == "Delivered" {
		result = config.DB.Exec("UPDATE payments SET paymentstatus = 'Successfull' WHERE id = ?", pid).Error
		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": result,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully Updated",
	})

}
