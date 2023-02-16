package controllers

import (
	"fmt"
	"path/filepath"

	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Addproducts(c *gin.Context) {
	var addproduct models.Product
	err := c.Bind(&addproduct)
	if err != nil { //Unmarshelling the Json Data
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Create(&addproduct)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":     "Successfully Added the Product",
		"Prodcuts ID": addproduct.Productid,
	})

}

type productdata struct {
	Productid   string
	Productname string
	Price       uint
	Baseimage   string
}

func ViewProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	var datas []productdata
	DB := config.DBconnect()
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
	result := DB.Raw(query).Scan(&datas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Products": datas,
	})
}

func AddImages(c *gin.Context) {
	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)
	pidconv := c.PostForm("pid")
	pid, _ := strconv.Atoi(pidconv)

	imagedata := models.Image{
		Image: image,
		Pid:   uint(pid),
	}
	DB := config.DBconnect()
	result := DB.Create(&imagedata)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Image Added Successfully",
	})

}

// "SELECT productname, description, stock, price, teams.teamname, sizes.sizetype, brands.brandname FROM products INNER JOIN teams ON products.teamid=teams.id INNER JOIN sizes ON products.sizeid=sizes.id INNER JOIN brands ON products.brandid=brands.id"
