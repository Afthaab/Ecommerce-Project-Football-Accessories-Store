package controllers

import (
	"fmt"

	"strconv"

	"net/http"

	"os"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func AddImages(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	// Check file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/jpg":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type",
		})
		return
	}
	// Create a new S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create S3 session",
		})
		return
	}
	// Create an S3 uploader
	uploader := s3manager.NewUploader(sess)

	// Open the file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not open file",
		})
		return
	}
	defer f.Close()

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(file.Filename),
		Body:   f,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
			"error": "Could not upload file to S3",
		})
		return
	}
	img := result.Location
	pid, _ := strconv.Atoi(c.Query("pid"))
	imagedata := models.Image{
		Pid:   uint(pid),
		Image: img,
	}
	result1 := config.DB.Create(&imagedata)
	if result1.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":   "Success",
		"Image ID":  imagedata.ID,
		"Image URL": imagedata.Image,
	})
}

func Addproducts(c *gin.Context) {
	var addproduct models.Product
	err := c.Bind(&addproduct)
	if err != nil { //Unmarshelling the Json Data
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.Create(&addproduct).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":     "Successfully Added the Product",
		"Prodcuts ID": addproduct.Productid,
	})

}

type productdata struct {
	Productid   string
	Productname string
	Price       uint
	Imageid     string
	Stock       uint
}

func ViewProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	var datas []productdata
	query := "SELECT * FROM products"
	if limit != 0 || offset != 0 {
		if limit == 0 {
			query = fmt.Sprintf("%s OFFSET %d", query, offset)
		} else if offset == 0 {
			query = fmt.Sprintf("%s LIMIT %d", query, limit)
		} else {
			query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
		}
	}
	result = config.DB.Raw(query).Scan(&datas).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Products": datas,
	})
}
