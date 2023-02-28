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

		//admin panel routes
		adminpanel := admin.Group("/adminpanel")
		{
			adminpanel.GET("/", middlewares.AdminAuth, controllers.Adminpanel)
			adminpanel.GET("/salesreport/downloadexcel", middlewares.AdminAuth, controllers.DownloadExcel)
			adminpanel.GET("/profilepage", middlewares.AdminAuth, controllers.AdminProfilepage)
		}

		//user management routes
		user := admin.Group("/user")
		{
			user.GET("/view", middlewares.AdminAuth, controllers.Adminviewuser)
			user.PUT("/block/unblock", middlewares.AdminAuth, controllers.UserManagement)
		}

		//specification management routes
		brands := admin.Group("/brands")
		{
			brands.POST("/add", middlewares.AdminAuth, controllers.AddBrands)
			brands.GET("/view", middlewares.AdminAuth, controllers.ViewSearchBrands)
			brands.PUT("/view/edit", middlewares.AdminAuth, controllers.Editbrands)
		}

		sizes := admin.Group("/sizes")
		{
			sizes.POST("/add", middlewares.AdminAuth, controllers.AddSize)
			sizes.GET("/view", middlewares.UserAuth, controllers.ViewSearchsize)
			sizes.PUT("/view/edit", middlewares.AdminAuth, controllers.Editsizes)

		}
		teams := admin.Group("/teams")
		{
			teams.POST("/add", middlewares.AdminAuth, controllers.AddTeams)
			teams.GET("/view", middlewares.AdminAuth, controllers.ViewSearchteams)
			teams.PUT("/view/edit", middlewares.AdminAuth, controllers.Editteams)

		}

		//product management
		products := admin.Group("/products")
		{
			products.POST("/add", middlewares.AdminAuth, controllers.Addproducts)
			products.POST("/add/images", middlewares.AdminAuth, controllers.AddImages)
		}

		//coupon routes
		coupon := admin.Group("/coupon")
		{
			coupon.POST("/add", middlewares.AdminAuth, controllers.AddCoupon)
			coupon.GET("/view", middlewares.AdminAuth, controllers.ViewCoupons)
			coupon.PUT("/view/edit", middlewares.AdminAuth, controllers.EditCoupon)
			coupon.DELETE("/view/delete", middlewares.AdminAuth, controllers.DeleteCoupon)
		}

		//order Management
		order := admin.Group("/order")
		{
			order.GET("/view/search", middlewares.AdminAuth, controllers.AdminOrderView)
			order.PUT("/cancel", middlewares.AdminAuth, controllers.CancelOrder)
			order.PUT("/statusupdate", middlewares.AdminAuth, controllers.StatusUpdate)
			order.PUT("/return", middlewares.AdminAuth, controllers.ReturnAccepted)
		}
	}
}
