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
		admin.GET("/user/view", middlewares.AdminAuth, controllers.Adminviewuser)
		admin.PUT("/user/block/unblock", middlewares.AdminAuth, controllers.UserManagement)

		//specification management routes
		admin.POST("/brands/add", middlewares.AdminAuth, controllers.AddBrands)
		admin.GET("/brands/view", middlewares.AdminAuth, controllers.ViewSearchBrands)
		admin.PUT("/brands/view/edit", middlewares.AdminAuth, controllers.Editbrands)

		admin.POST("/sizes/add", middlewares.AdminAuth, controllers.AddSize)
		admin.GET("/sizes/view", middlewares.AdminAuth, controllers.ViewSearchsize)
		admin.PUT("/sizes/view/edit", middlewares.AdminAuth, controllers.Editsizes)

		admin.POST("/teams/add", middlewares.AdminAuth, controllers.AddTeams)
		admin.GET("/teams/view", middlewares.AdminAuth, controllers.ViewSearchteams)
		admin.PUT("/teams/view/edit", middlewares.AdminAuth, controllers.Editteams)

		//product management
		admin.POST("/products/add", middlewares.AdminAuth, controllers.Addproducts)
		admin.POST("/products/add/images", middlewares.AdminAuth, controllers.AddImages)

		//coupon routes
		admin.POST("/coupon/add", middlewares.AdminAuth, controllers.AddCoupon)
		admin.GET("/coupon/view", middlewares.AdminAuth, controllers.ViewCoupons)
		admin.PUT("/coupon/view/edit", middlewares.AdminAuth, controllers.EditCoupon)
		admin.DELETE("/coupon/view/delete", middlewares.AdminAuth, controllers.DeleteCoupon)

		//order Management
		admin.GET("/order/view/search", middlewares.AdminAuth, controllers.AdminOrderView)
		admin.PUT("/order/cancel", middlewares.AdminAuth, controllers.CancelOrder)
		admin.PUT("/order/statusupdate", middlewares.AdminAuth, controllers.StatusUpdate)
		admin.PUT("/order/return", middlewares.AdminAuth, controllers.ReturnAccepted)

		//excel
		admin.GET("/order/salesreport/download/excel", controllers.DownloadExcel)
		admin.GET("/order/salesreport/download/pdf", controllers.DownloadPdf)
	}
}
