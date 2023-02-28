package controllers

import (
	"net/http"
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
)

func AddBrands(c *gin.Context) {
	var addbrand models.Brand
	if c.Bind(&addbrand) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Could not bind JSON data",
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.Create(&addbrand).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":    "New Brand added Successfully",
		"Brand ID":   addbrand.ID,
		"Brand Name": addbrand.Brandname,
	})
}

func AddSize(c *gin.Context) {
	var addsize models.Size
	if c.Bind(&addsize) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in binding JSON data",
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.Create(&addsize).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":   "New Size added Successfully",
		"Size ID":   addsize.ID,
		"Size Type": addsize.Sizetype,
	})

}

func AddTeams(c *gin.Context) {
	var addteam models.Team
	if c.Bind(&addteam) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.Create(&addteam).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":   "New Team added Successfully",
		"Team ID":   addteam.ID,
		"Team Name": addteam.Teamname,
	})
}

func ViewSearchBrands(c *gin.Context) {
	searchid := c.Query("brandid")
	var viewbrands []models.Brand
	// DB := config.DBconnect()
	result = config.DB.Find(&viewbrands).Error
	if searchid != "" {
		result = config.DB.First(&viewbrands, "id = ?", searchid).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Brand": viewbrands,
		})
		return
	}
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Brands": viewbrands,
	})

}

func ViewSearchsize(c *gin.Context) {
	searchid := c.Query("sizeid")
	var viewsizes []models.Size
	// DB := config.DBconnect()
	result = config.DB.Find(&viewsizes).Error
	if searchid != "" {
		result = config.DB.First(&viewsizes, "id = ?", searchid).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Size": viewsizes,
		})
		return
	}
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Sizes": viewsizes,
	})

}
func ViewSearchteams(c *gin.Context) {
	searchid := c.Query("teamid")
	var viewteams []models.Team
	// DB := config.DBconnect()
	result = config.DB.Find(&viewteams).Error
	if searchid != "" {
		result = config.DB.First(&viewteams, "id = ?", searchid).Error
		if result != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": result,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Team": viewteams,
		})
		return
	}
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Teams": viewteams,
	})

}
func Editbrands(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("brandid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editbrands models.Brand
	if c.Bind(&editbrands) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editbrands.ID = uint(id)
	// DB := config.DBconnect()
	result = config.DB.Model(&editbrands).Updates(models.Brand{Brandname: editbrands.Brandname}).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully updated the Brand",
	})

}
func Editsizes(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("sizeid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editsizes models.Size
	if c.Bind(&editsizes) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editsizes.ID = uint(id)
	// DB := config.DBconnect()
	result = config.DB.Model(&editsizes).Updates(models.Size{Sizetype: editsizes.Sizetype}).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully updated the Size",
	})

}
func Editteams(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("teamid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editteams models.Team
	if c.Bind(&editteams) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editteams.ID = uint(id)
	// DB := config.DBconnect()
	result = config.DB.Model(&editteams).Updates(models.Team{Teamname: editteams.Teamname}).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully updated the Team",
	})
}
