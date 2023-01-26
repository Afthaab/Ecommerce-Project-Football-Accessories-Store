package routes

import (
	"github.com/afthab/e_commerce/controllers"
	"github.com/gin-gonic/gin"
)

func OpenRoutes(r *gin.Engine) {
	r.POST("/viewallproducts", controllers.ViewProducts)
	r.POST("/viewallproducts/contraints", controllers.ViewAllProductConstraints)
}