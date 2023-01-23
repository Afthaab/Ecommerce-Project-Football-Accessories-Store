package controllers

import "github.com/gin-gonic/gin"

func Addproducts(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "This is Product adding page",
	})
}
