package routes

import (
	"github.com/afthab/e_commerce/controllers"
	"github.com/afthab/e_commerce/middlewares"
	"github.com/gin-gonic/gin"
)

func Adminroutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.POST("/signin", controllers.Adminsignin)
		admin.GET("/signout", middlewares.AdminAuth, controllers.AdminSignout)

		//routes with middlewares
		admin.GET("/adminpanel", middlewares.AdminAuth, controllers.Adminpanel)
		admin.GET("/proflepage", middlewares.AdminAuth, controllers.AdminProfilepage)
		admin.GET("/viewuser", middlewares.AdminAuth, controllers.Adminviewuser)
		admin.POST("/searchuser", middlewares.AdminAuth, controllers.Adminsearchuser)
		admin.PUT("/searchuser/blockuser", middlewares.AdminAuth, controllers.Adminblockuser)
		admin.PUT("/searchuser/unblockuser", middlewares.AdminAuth, controllers.Adminunblockuser)
		admin.POST("/addproducts", middlewares.AdminAuth, controllers.Addproducts)
	}
}
