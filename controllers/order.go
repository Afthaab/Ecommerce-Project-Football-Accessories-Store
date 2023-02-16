package controllers

import (
	"fmt"
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
	DB := config.DBconnect()

	//getting the payment details through the payment id recieved in the query params
	var paymentdata models.Payment
	result := DB.Raw("SELECT * FROM payments WHERE id = ?", pid).Scan(&paymentdata).Error
	if result != nil {
		c.JSON(404, gin.H{
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
		result4 := DB.Raw("select addressid from addresses where defaultadd=true and uid = ?", id).Scan(&orderdata.Addid).Error
		if result4 != nil {
			c.JSON(404, gin.H{
				"Error": result4.Error(),
			})
			return
		}
	}
	//creating the order
	result1 := DB.Create(&orderdata).Error
	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": result1,
		})
		return
	}
	//shifting all the products in the carts to the orderditems table to save the products from the cart
	var itemdata models.Orderditems
	result2 := DB.Raw("insert into orderditems(pid, oid, quantity, price, totalprice, uid) select carts.pid, orders.orderid, carts.quantity, carts.price, carts.totalprice, orders.useridno from orders inner join carts on orders.useridno=carts.cartid where orderid = ?", orderdata.Orderid).Scan(&itemdata).Error
	if result2 != nil {
		c.JSON(400, gin.H{
			"Error": result2,
		})
		return
	}
	//deleting the products in the cart
	result3 := DB.Exec("DELETE FROM carts WHERE cartid = ?", id).Error
	if result3 != nil {
		c.JSON(400, gin.H{
			"Error": result3,
		})
		return
	}
	c.JSON(200, gin.H{
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
	DB := config.DBconnect()

	//to check if the oid or pid is passed through the query params
	if oid != uuid.Nil || prodid != 0 {

		//showing the shipment address
		var addressesdata models.Address
		result2 := DB.Raw("select addresses.* from orders inner join addresses on addresses.addressid=orders.addid where orders.orderid = ?", oid).Scan(&addressesdata).Error
		if result2 != nil {
			c.JSON(404, gin.H{
				"Error": result2,
			})
			return
		}
		c.JSON(200, gin.H{
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
		result3 := DB.Raw(query).Scan(&data).Error
		if result3 != nil {
			c.JSON(404, gin.H{
				"Error": result3,
			})
			return
		}
		c.JSON(200, gin.H{
			"Order Details": data,
		})

		//showing the payment details for the particular order
		var paymentdata models.Payment
		result4 := DB.Raw("select payments.* from orders INNER JOIN payments on orders.paymentid=payments.id where orderid = ?", oid).Scan(&paymentdata).Error
		if result4 != nil {
			c.JSON(400, gin.H{
				"Error": result4,
			})
			return
		}
		c.JSON(200, gin.H{
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
		result := DB.Raw("select products.productname, orderditems.*,orders.orderstatus, orders.useridno, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid where uid = ?", uid).Scan(&data).Error
		if result != nil {
			c.JSON(404, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(200, gin.H{
			"Your Orders": data,
		})
	}

}

func CancelOrder(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	DB := config.DBconnect()

	//changing the status in the orders table
	result := DB.Exec("UPDATE orders SET orderstatus = 'Order Cancelled' WHERE orderid = ?", oid)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//to get the payment id and method to update the status in the payment table
	var paymentdata models.Payment
	result2 := DB.Raw("select payments.paymentmethod,id,payments.useridno from orders inner join payments on payments.id=orders.paymentid where orderid = ?", oid).Scan(&paymentdata)
	if result2.Error != nil {
		c.JSON(400, gin.H{
			"Error": result2.Error.Error(),
		})
		return
	}
	fmt.Println(paymentdata.Useridno)
	//if the payment in done in razor pay inserting the balance to the wallet and updating the status in the payment table
	if paymentdata.Paymentmethod == "Razor Pay" {
		result1 := DB.Exec("INSERT INTO wallets (balance, oid, walletid) SELECT SUM(totalprice), ?, ? FROM orderditems WHERE oid = ?", oid, paymentdata.Useridno, oid)
		if result1.Error != nil {
			c.JSON(400, gin.H{
				"Error": result1.Error.Error(),
			})
			return
		}
		result3 := DB.Exec("UPDATE payments set paymentstatus = 'refunded' WHERE id = ?", paymentdata.ID)
		if result3.Error != nil {
			c.JSON(400, gin.H{
				"Error": result3.Error.Error(),
			})
			return
		}
	} else {
		//if the payment is done in cod just updating the payment status
		result3 := DB.Exec("UPDATE payments set paymentstatus = 'Failed' WHERE id = ?", paymentdata.ID)
		if result3.Error != nil {
			c.JSON(400, gin.H{
				"Error": result3.Error.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"Message":   "Success",
		"Orderd ID": oid,
	})
}

func ReturnOrder(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	uid, _ := strconv.Atoi(c.GetString("userid"))
	DB := config.DBconnect()
	result := DB.Exec("UPDATE orders set orderstatus = 'Return requested' WHERE orderid = ? AND useridno = ?", oid, uid)
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

func ShowWallet(c *gin.Context) {
	type data struct {
		Id        string
		Balance   uint
		Walletid  uint
		Oid       uuid.UUID
		CreatedAt time.Time
	}
	var walletdata []data
	uid, _ := strconv.Atoi(c.GetString("userid"))
	DB := config.DBconnect()
	//Getting the wallet data
	result := DB.Raw("select * from wallets where walletid = ?", uid).Scan(&walletdata)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	//caluculating the total sum
	var totalbalance uint
	result1 := DB.Raw("select SUM(balance) from wallets where walletid = ?", uid).Scan(&totalbalance)
	if result1.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Wallet":        walletdata,
		"Total Balance": totalbalance,
	})
}

// ------------------------ Admin Route Functions ----------------------
func AdminOrderView(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	pid, _ := strconv.Atoi(c.Query("pid"))

	var data []Orderdetails
	DB := config.DBconnect()
	// to search a particular order
	if oid != uuid.Nil || pid != 0 {
		var addressesdata models.Address
		// getting the address from orders
		result2 := DB.Raw("select addresses.* from orders inner join addresses on addresses.addressid=orders.addid where orders.orderid = ?", oid).Scan(&addressesdata).Error
		if result2 != nil {
			c.JSON(404, gin.H{
				"Error": result2,
			})
			return
		}
		c.JSON(200, gin.H{
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
		result3 := DB.Raw(query).Scan(&data).Error
		if result3 != nil {
			c.JSON(404, gin.H{
				"Error": result3,
			})
			return
		}
		//showing the payment details for the particular order
		var paymentdata models.Payment
		result4 := DB.Raw("select payments.* from orders INNER JOIN payments on orders.paymentid=payments.id where orderid = ?", oid).Scan(&paymentdata).Error
		if result4 != nil {
			c.JSON(400, gin.H{
				"Error": result4,
			})
			return
		}
		c.JSON(200, gin.H{
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
		result := DB.Raw("select products.productname, orderditems.*,orders.orderstatus,orders.useridno, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid;").Scan(&data)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"Orders": data,
	})

}
func ReturnAccepted(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	DB := config.DBconnect()
	var orderdata models.Orders
	result := DB.Raw("UPDATE orders SET orderstatus = 'Returned' WHERE orderid = ? RETURNING totalamount,useridno, paymentid", oid).Scan(&orderdata)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	result1 := DB.Exec("INSERT INTO wallets(balance, oid, walletid) VALUES(?, ?, ?)", orderdata.Totalamount, oid, orderdata.Useridno)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	fmt.Println("hii")
	result2 := DB.Exec("UPDATE payments SET paymentstatus = 'Refunded' WHERE id = ?", orderdata.Paymentid)
	if result2.Error != nil {
		c.JSON(400, gin.H{
			"Error": result2.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Success",
	})

}
func StatusUpdate(c *gin.Context) {
	status := c.Query("status")
	oid, _ := uuid.FromString(c.Query("orderid"))
	DB := config.DBconnect()
	var pid uint
	result := DB.Raw("UPDATE orders SET orderstatus = ? WHERE orderid = ? RETURNING paymentid", status, oid).Scan(&pid)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if status == "delivered" {
		result1 := DB.Exec("UPDATE payments SET paymentstatus = 'successfull' WHERE id = ?", pid)
		if result1.Error != nil {
			c.JSON(400, gin.H{
				"Error": result1.Error.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"Message": "Successfully Updated",
	})

}
