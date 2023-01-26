package controllers

import (
	"path/filepath"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Addproducts(c *gin.Context) {
	productname := c.PostForm("product_name")
	description := c.PostForm("description")
	stockconv := c.PostForm("stock")
	stock, _ := strconv.Atoi(stockconv)
	team := c.PostForm("team_name")
	brand := c.PostForm("brand_name")
	size := c.PostForm("size_type")
	priceconv := c.PostForm("price")
	price, _ := strconv.Atoi(priceconv)

	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)

	coverpath, _ := c.FormFile("cover")
	extension1 := filepath.Ext(coverpath.Filename)
	cover := uuid.New().String() + extension1
	c.SaveUploadedFile(coverpath, "./public/images"+cover)

	var sizedata models.Size
	DB := config.DBconnect()
	result := DB.Raw("SELECT sizeid FROM sizes WHERE sizetype = ?", size).Scan(&sizedata)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error":   "Size is not found please verify",
			"Message": "If its a new size add the size first",
		})
		return
	}

	var teamdata models.Team
	result1 := DB.Raw("SELECT teamid FROM teams WHERE teamname = ?", team).Scan(&teamdata)
	if result1.Error != nil {
		c.JSON(404, gin.H{
			"Error":   "Team is not found please verify",
			"Message": "If its a new Team add the Team first",
		})
		return
	}

	var brandata models.Brand
	result2 := DB.Raw("SELECT brandid FROM brands WHERE brandname = ?", brand).Scan(&brandata)
	if result2.Error != nil {
		c.JSON(404, gin.H{
			"Error":   "Team is not found please verify",
			"Message": "If its a new Team add the Team first",
		})
		return
	}
	addproduct := models.Product{
		Productname: productname,
		Description: description,
		Stock:       uint(stock),
		Price:       uint(price),
		Image:       image,
		Cover:       cover,
		Teamid:      teamdata.Teamid,
		Brandid:     brandata.Brandid,
		Sizeid:      sizedata.Sizeid,
	}
	result3 := DB.Create(&addproduct)
	if result3.Error != nil {
		c.JSON(500, gin.H{
			"Error": result3.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Product had added successfully",
	})

}

func AddBrands(c *gin.Context) {
	name := c.PostForm("brand_name")
	addbrand := models.Brand{
		Brandname: name,
	}
	DB := config.DBconnect()
	result := DB.Create(&addbrand)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "New Brand added Successfully",
	})
}

func AddSize(c *gin.Context) {
	name := c.PostForm("size_type")
	addsize := models.Size{
		Sizetype: name,
	}
	DB := config.DBconnect()
	result := DB.Create(&addsize)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "New Size added Successfully",
	})

}

func AddTeams(c *gin.Context) {
	name := c.PostForm("team_name")
	addteam := models.Team{
		Teamname: name,
	}
	DB := config.DBconnect()
	result := DB.Create(&addteam)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "New Team added Successfully",
	})
}

func ViewProducts(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "View Products Page",
	})
}

func ViewAllProductConstraints(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "View all products constraints page",
	})
}
