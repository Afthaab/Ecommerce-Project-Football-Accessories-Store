package routes

import (
	"github.com/afthab/e_commerce/controllers"
	"github.com/afthab/e_commerce/middlewares"
	"github.com/gin-gonic/gin"
)

func Adminroutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		//admin routes
		admin.POST("/signin", controllers.Adminsignin)
		admin.GET("/signout", middlewares.AdminAuth, controllers.AdminSignout)
		admin.GET("/adminpanel", middlewares.AdminAuth, controllers.Adminpanel)
		admin.GET("/profilepage", middlewares.AdminAuth, controllers.AdminProfilepage)

		//user management routes
		admin.GET("/viewuser", middlewares.AdminAuth, controllers.Adminviewuser)
		admin.GET("/searchuser", middlewares.AdminAuth, controllers.Adminsearchuser)
		admin.PUT("/searchuser/blockuser", middlewares.AdminAuth, controllers.Adminblockuser)
		admin.PUT("/searchuser/unblockuser", middlewares.AdminAuth, controllers.Adminunblockuser)

		//specification management routes
		admin.POST("/addbrands", middlewares.AdminAuth, controllers.AddBrands)
		admin.POST("/addsize", middlewares.AdminAuth, controllers.AddSize)
		admin.POST("/addteams", middlewares.AdminAuth, controllers.AddTeams)
		admin.GET("/view/searchbrands", middlewares.AdminAuth, controllers.ViewSearchBrands)
		admin.GET("/view/searchsizes", middlewares.AdminAuth, controllers.ViewSearchsize)
		admin.GET("/view/searchteams", middlewares.AdminAuth, controllers.ViewSearchteams)
		admin.PUT("/editbrands", middlewares.AdminAuth, controllers.Editbrands)
		admin.PUT("/editsizes", middlewares.AdminAuth, controllers.Editsizes)
		admin.PUT("/editteams", middlewares.AdminAuth, controllers.Editteams)

		//product management
		admin.POST("/addproducts", middlewares.AdminAuth, controllers.Addproducts)
		admin.POST("/addimages", middlewares.AdminAuth, controllers.AddImages)

		//payment routes

	}
}
