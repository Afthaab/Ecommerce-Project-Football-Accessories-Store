package controllers

import (
	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
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
		"Message": "Successfully Added the Product",
	})

}

func ViewProducts(c *gin.Context) {
	var viewproducts []models.Product
	DB := config.DBconnect()
	result := DB.Raw("SELECT productname,description,stock,price,brands.brands,teams.teams,sizes.sizes FROM products INNER JOIN brands ON products.brandid=brands.brandid INNER JOIN teams ON products.teamid=teams.teamid INNER JOIN sizes ON products.sizeid=sizes.sizeid").Scan(&viewproducts)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Products": viewproducts,
	})
}

func ViewAllProductConstraints(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "View all products constraints page",
	})
}
