package controllers

import (
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func AddBrands(c *gin.Context) {
	var addbrand models.Brand
	if c.Bind(&addbrand) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind JSON data",
		})
		return
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
	var addsize models.Size
	if c.Bind(&addsize) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding JSON data",
		})
		return
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
	var addteam models.Team
	if c.Bind(&addteam) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
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

func ViewSearchBrands(c *gin.Context) {
	searchid := c.Query("brandid")
	var viewbrands []models.Brand
	DB := config.DBconnect()
	result := DB.Find(&viewbrands)
	if searchid != "" {
		result1 := DB.First(&viewbrands, "id = ?", searchid)
		if result1.Error != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Brand": viewbrands,
		})
		return
	}
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Brands": viewbrands,
	})

}

func ViewSearchsize(c *gin.Context) {
	searchid := c.Query("sizeid")
	var viewsizes []models.Size
	DB := config.DBconnect()
	result := DB.Find(&viewsizes)
	if searchid != "" {
		result1 := DB.First(&viewsizes, "id = ?", searchid)
		if result1.Error != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Size": viewsizes,
		})
		return
	}
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Sizes": viewsizes,
	})

}
func ViewSearchteams(c *gin.Context) {
	searchid := c.Query("teamid")
	var viewteams []models.Team
	DB := config.DBconnect()
	result := DB.Find(&viewteams)
	if searchid != "" {
		result1 := DB.First(&viewteams, "id = ?", searchid)
		if result1.Error != nil {
			c.JSON(404, gin.H{
				"Error": result1.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Team": viewteams,
		})
		return
	}
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Teams": viewteams,
	})

}
func Editbrands(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("brandid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editbrands models.Brand
	if c.Bind(&editbrands) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editbrands.ID = uint(id)
	DB := config.DBconnect()
	result := DB.Model(&editbrands).Updates(models.Brand{Brandname: editbrands.Brandname})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully updated the Brand",
	})

}
func Editsizes(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("sizeid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editsizes models.Size
	if c.Bind(&editsizes) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editsizes.ID = uint(id)
	DB := config.DBconnect()
	result := DB.Model(&editsizes).Updates(models.Size{Sizetype: editsizes.Sizetype})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully updated the Size",
	})

}
func Editteams(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("teamid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editteams models.Team
	if c.Bind(&editteams) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editteams.ID = uint(id)
	DB := config.DBconnect()
	result := DB.Model(&editteams).Updates(models.Team{Teamname: editteams.Teamname})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully updated the Team",
	})
}
