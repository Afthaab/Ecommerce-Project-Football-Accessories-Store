package controllers

import (
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
	var paymentdata models.Payment
	result := DB.Raw("SELECT * FROM payments WHERE id = ?", pid).Scan(&paymentdata).Error
	if result != nil {
		c.JSON(404, gin.H{
			"Error": result,
		})
		return
	}
	orderdata := models.Orders{
		Useridno:    uint(id),
		Totalamount: paymentdata.Totalamount,
		Paymentid:   uint(pid),
		Orderstatus: "order placed",
	}
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
	result1 := DB.Create(&orderdata).Error
	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": result1,
		})
		return
	}
	c.JSON(200, gin.H{
		"Oid": orderdata.Orderid,
	})
	var itemdata models.Orderditems
	result2 := DB.Raw("insert into orderditems(pid, oid, quantity, price, totalprice, uid) select carts.pid, orders.orderid, carts.quantity, carts.price, carts.totalprice, orders.useridno from orders inner join carts on orders.useridno=carts.cartid where orderid = ?", orderdata.Orderid).Scan(&itemdata).Error
	if result2 != nil {
		c.JSON(400, gin.H{
			"Error": result2,
		})
		return
	}
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

func ViewOrder(c *gin.Context) {
	oid, _ := uuid.FromString(c.Query("orderid"))
	uid, _ := strconv.Atoi(c.GetString("userid"))
	type orderdetails struct {
		Productname string
		Id          uint
		Pid         uint
		Oid         uuid.UUID
		Quantity    uint
		Price       uint
		Totalprice  uint
		Orderstatus string
		Date        time.Time
	}
	var data []orderdetails
	DB := config.DBconnect()
	if oid != uuid.Nil {
		var addressesdata models.Address
		result2 := DB.Raw("select addresses.* from orders inner join addresses on addresses.addressid=orders.addid where orders.orderid = ?", oid).Scan(&addressesdata).Error
		if result2 != nil {
			c.JSON(404, gin.H{
				"Error": result2,
			})
			return
		}
		c.JSON(200, gin.H{
			"Shipping Address": addressesdata,
		})
		result3 := DB.Raw("select products.productname, orderditems.*, orders.orderstatus, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid where oid = ?", oid).Scan(&data).Error
		if result3 != nil {
			c.JSON(404, gin.H{
				"Error": result3,
			})
			return
		}
		c.JSON(200, gin.H{
			"Order Details": data,
		})
		var paymentdata models.Payment
		result4 := DB.Raw("select payments.* from orders INNER JOIN payments on orders.paymentid=payments.id where orderid = ?", oid).Scan(&paymentdata).Error
		if result4 != nil {
			c.JSON(400, gin.H{
				"Error": result4,
			})
			return
		}

		c.JSON(200, gin.H{
			"Payment Details": paymentdata,
		})

	} else {
		result := DB.Raw("select products.productname, orderditems.*, orders.orderstatus, date(created_at) from orderditems inner join products on orderditems.pid = products.productid inner join orders on orderditems.oid=orders.orderid where uid = ?", uid).Scan(&data).Error
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
	uid, _ := strconv.Atoi(c.GetString("userid"))
	DB := config.DBconnect()
	result := DB.Exec("update orders set orderstatus = 'order cancelled' where orderid = ? And useridno = ?", oid, uid)
	if result.RowsAffected == 0 {
		c.JSON(400, gin.H{
			"Error": "No rows affected",
		})
		return
	}
	
	// result2 := DB.Exec("UPDATE payments SET paymentstatus = 'failed' WHERE id IN (SELECT paymentid FROM orders WHERE orderid = ? AND useridno = ?)", oid, uid)
	// if result2.RowsAffected == 0 {
	// 	c.JSON(400, gin.H{
	// 		"Error": "No rows affected",
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"Message":  "Successfully Updated order and Payment ",
		"Order Id": oid,
	})

}
