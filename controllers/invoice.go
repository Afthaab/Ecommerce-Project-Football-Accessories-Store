package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"text/template"

	"fmt"
	"os/exec"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
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

type Orderdata struct {
	Userid        uint
	Firstname     string
	Lastname      string
	Email         string
	Phone         string
	Productid     string
	Productname   string
	Quantity      uint
	Price         uint
	Totalprice    uint
	Oid           uuid.UUID
	Id            uint
	Paymentmethod string
	Paymentstatus string
}

func DownloadExcel(c *gin.Context) {
	var data []Orderdata
	DB := config.DBconnect()
	result1 := DB.Raw("select users.userid, users.firstname, users.lastname, users.email, users.phone, products.productid, products.productname, orderditems.quantity, orderditems.price, orderditems.totalprice, orderditems.oid, payments.id, payments.paymentmethod, payments.paymentstatus from orderditems inner join users on orderditems.uid=users.userid inner join products on products.productid=orderditems.pid inner join orders on orders.orderid=orderditems.oid inner join payments on orders.paymentid=payments.id where payments.paymentstatus = 'successfull'").Scan(&data)
	if result1.Error != nil {
		c.JSON(404, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	file := excelize.NewFile()
	sheetName := "Sales Report"
	index := file.NewSheet(sheetName)

	// Set some cell values
	file.SetCellValue(sheetName, "A1", "Product ID")
	file.SetCellValue(sheetName, "B1", "Product Name")
	file.SetCellValue(sheetName, "C1", "Quantity")
	file.SetCellValue(sheetName, "D1", "Price")
	file.SetCellValue(sheetName, "E1", "Total Price")
	file.SetCellValue(sheetName, "F1", "Order ID")
	file.SetCellValue(sheetName, "G1", "Payment ID")
	file.SetCellValue(sheetName, "H1", "Payment Method")
	file.SetCellValue(sheetName, "I1", "Payment Status")
	file.SetCellValue(sheetName, "J1", "User ID")
	file.SetCellValue(sheetName, "K1", "First Name")
	file.SetCellValue(sheetName, "L1", "Last Name")
	file.SetCellValue(sheetName, "M1", "Email")
	file.SetCellValue(sheetName, "N1", "Phone")

	for i, person := range data {
		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), person.Productid)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), person.Productname)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), person.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), person.Price)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), person.Totalprice)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), person.Oid)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), person.Id)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), person.Paymentmethod)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), person.Paymentstatus)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), person.Userid)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", row), person.Firstname)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", row), person.Lastname)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", row), person.Email)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", row), person.Phone)
	}

	// Set the active sheet
	file.SetActiveSheet(index)

	// Set the content type and headers for the response
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=example.xlsx")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Write the Excel file to the response
	err := file.Write(c.Writer)
	if err != nil {
		fmt.Println(err)
	}
}
func DownloadPdf(c *gin.Context) {
	var data []Orderdata
	DB := config.DBconnect()
	result1 := DB.Raw("select users.userid, users.firstname, users.lastname, users.email, users.phone, products.productid, products.productname, orderditems.quantity, orderditems.price, orderditems.totalprice, orderditems.oid, payments.id, payments.paymentmethod, payments.paymentstatus from orderditems inner join users on orderditems.uid=users.userid inner join products on products.productid=orderditems.pid inner join orders on orders.orderid=orderditems.oid inner join payments on orders.paymentid=payments.id where payments.paymentstatus = 'successfull'").Scan(&data)
	if result1.Error != nil {
		c.JSON(404, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	file := excelize.NewFile()
	sheetName := "Sales Report"

	// Set some cell values
	file.SetCellValue(sheetName, "A1", "Product ID")
	file.SetCellValue(sheetName, "B1", "Product Name")
	file.SetCellValue(sheetName, "C1", "Quantity")
	file.SetCellValue(sheetName, "D1", "Price")
	file.SetCellValue(sheetName, "E1", "Total Price")
	file.SetCellValue(sheetName, "F1", "Order ID")
	file.SetCellValue(sheetName, "G1", "Payment ID")
	file.SetCellValue(sheetName, "H1", "Payment Method")
	file.SetCellValue(sheetName, "I1", "Payment Status")
	file.SetCellValue(sheetName, "J1", "User ID")
	file.SetCellValue(sheetName, "K1", "First Name")
	file.SetCellValue(sheetName, "L1", "Last Name")
	file.SetCellValue(sheetName, "M1", "Email")
	file.SetCellValue(sheetName, "N1", "Phone")

	for i, person := range data {
		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), person.Productid)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), person.Productname)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), person.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), person.Price)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), person.Totalprice)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), person.Oid)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), person.Id)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), person.Paymentmethod)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), person.Paymentstatus)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), person.Userid)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", row), person.Firstname)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", row), person.Lastname)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", row), person.Email)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", row), person.Phone)
	}

	// Save the Excel file to a temporary file
	tempFile, err := ioutil.TempFile("", "excel-*.xlsx")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer os.Remove(tempFile.Name())
	err = file.SaveAs(tempFile.Name())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Convert the Excel file to PDF using unoconv
	cmd := exec.Command("unoconv", "-f", "pdf", tempFile.Name())
	pdfBytes, err := cmd.Output()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Serve the PDF file as a download
	c.Header("Content-Disposition", "attachment; filename=example.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)

}
