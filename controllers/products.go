package controllers

import (
	"net/http"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func Addproducts(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "This is Product adding page",
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
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "New Team added Successfully",
	})
}
