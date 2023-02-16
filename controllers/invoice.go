package controllers

import (
	"bytes"
	"strconv"

	"text/template"

	"os/exec"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Invoice struct {
	First_Name   string
	Last_Name    string
	Email        string
	Phone        string
	Address      Addressdata
	Ordersdata   Ordersdata
	Paymentsdata Paymentsdata
	Products     []Productsdata
}

type Addressdata struct {
	Name     string
	Phoneno  string
	Houseno  string
	Area     string
	Landmark string
	City     string
	Pincode  string
	District string
	State    string
	Country  string
}

type Productsdata struct {
	Product_Name string
	Quantity     uint
	Price        uint
}

type Ordersdata struct {
	Order_Id     uuid.UUID
	Order_Status string
	Date         string
}

type Paymentsdata struct {
	Total_Amount   uint
	Payment_Method string
	Payment_Status string
}

const invoiceTemplate = `
Order Details:<br>
Order ID : {{.Ordersdata.Order_Id}}<br>
Order Date:{{.Ordersdata.Date}} <br>
Order Status:{{.Ordersdata.Order_Status}}
<hr>
User Details:<br>
First Name : {{.First_Name}} <br>
Last Name : {{.Last_Name}}<br>
Email : {{.Email}}<br>
Phone :{{.Phone}}<br>
<hr>
Billing Address :
Name : {{.Address.Name}}
Phone number : {{.Address.Phoneno}} <br>
House number : {{.Address.Houseno}} <br>
Area : {{.Address.Area}} <br>
Landmark : {{.Address.Landmark}} <br>
City : {{.Address.City}} <br>
Pincode : {{.Address.Pincode}} <br>
District : {{.Address.District}} <br>
State : {{.Address.State}} <br>
Country : {{.Address.Country}} <br>
<hr>
Product Details:
{{range .Products}}
Product Name :{{.Product_Name}} <br>
Price : {{.Price}}<br>
Quantity : {{.Quantity}}<br>
{{end}}
<hr>
Payment Details:
Total Amount : {{.Paymentsdata.Total_Amount}}
Payment method : {{.Paymentsdata.Payment_Method}}<br>
Payment Status : {{.Paymentsdata.Payment_Status}}<br>
<hr>
`

func GenerateInvoice(c *gin.Context) {
	uid, _ := strconv.Atoi(c.GetString("userid"))
	oid, _ := uuid.FromString(c.Query("orderid"))
	var userdatas models.User
	var addressdatas models.Address
	var orderdatas models.Orders
	var billingdata models.Payment
	var itemsdata []models.Product
	DB := config.DBconnect()
	result := DB.Raw("SELECT * FROM users WHERE userid = ?", uid).Scan(&userdatas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	result1 := DB.Raw("SELECT * FROM orders WHERE orderid = ?", oid).Scan(&orderdatas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}

	result2 := DB.Raw("SELECT * FROM addresses WHERE addressid = ?", orderdatas.Addid).Scan(&addressdatas)
	if result2.Error != nil {
		c.JSON(404, gin.H{
			"Error": result2.Error.Error(),
		})
		return
	}

	result3 := DB.Raw("SELECT * FROM payments WHERE id = ?", orderdatas.Paymentid).Scan(&billingdata)
	if result3.Error != nil {
		c.JSON(404, gin.H{
			"Error": result3.Error.Error(),
		})
		return
	}
	result4 := DB.Raw("SELECT products.* from orderditems inner join products on orderditems.pid=products.productid  WHERE oid = ?", oid).Scan(&itemsdata)
	if result4.Error != nil {
		c.JSON(404, gin.H{
			"Error": result4.Error.Error(),
		})
		return
	}
	products := make([]Productsdata, len(itemsdata))
	for i, data := range itemsdata {
		products[i] = Productsdata{
			Product_Name: data.Productname,
			Quantity:     data.Stock,
			Price:        data.Price,
		}
	}
	datestring := orderdatas.CreatedAt.Format("2006-01-02")
	invoice := Invoice{
		First_Name: userdatas.Firstname,
		Last_Name:  userdatas.Lastname,
		Email:      userdatas.Email,
		Phone:      userdatas.Phone,
		Address: Addressdata{
			Name:     addressdatas.Name,
			Phoneno:  addressdatas.Phoneno,
			Houseno:  addressdatas.Houseno,
			Area:     addressdatas.Area,
			Landmark: addressdatas.Landmark,
			City:     addressdatas.City,
			Pincode:  addressdatas.Pincode,
			District: addressdatas.District,
			State:    addressdatas.State,
			Country:  addressdatas.Country,
		},
		Ordersdata: Ordersdata{
			Order_Id:     uuid.UUID(orderdatas.Orderid),
			Order_Status: orderdatas.Orderstatus,
			Date:         datestring,
		},
		Paymentsdata: Paymentsdata{
			Total_Amount:   billingdata.Totalamount,
			Payment_Method: billingdata.Paymentmethod,
			Payment_Status: billingdata.Paymentstatus,
		},
		Products: products,
	}
	templatepage, err := template.New("invoice").Parse(invoiceTemplate)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	var buf bytes.Buffer
	err = templatepage.Execute(&buf, invoice)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}
	cmd := exec.Command("wkhtmltopdf", "-", "invoice.pdf")
	cmd.Stdin = &buf
	err = cmd.Run()
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}
	c.HTML(200, "invoice.html", gin.H{})
}
func InvoiceDownload(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("invoice.pdf")
}
