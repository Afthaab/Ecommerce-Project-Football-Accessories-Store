package routes

import (
	"github.com/afthab/e_commerce/controllers"
	"github.com/gin-gonic/gin"
)

func Adminroutes(r *gin.Engine) {
	r.POST("/admin/signin", controllers.Adminsignin)
	r.GET("/admin/adminpanel", controllers.Adminpanel)
	r.GET("/admin/viewuser", controllers.Adminviewuser)
	r.POST("/admin/searchuser", controllers.Adminsearchuser)
	r.PUT("/admin/blockuser", controllers.Adminblockuser)
	r.PUT("/admin/unblockuser", controllers.Adminunblockuser)
}
