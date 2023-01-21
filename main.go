package main

import (
	"os"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/controllers"
	"github.com/afthab/e_commerce/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	config.DBconnect()
}
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//user
	r.POST("/user/signup", controllers.Usersignup)
	r.POST("/user/signup/otpvalidate", controllers.Otpvalidate)
	r.POST("/user/signin", controllers.Usersignin)
	//admin
	r.POST("/admin/signin", controllers.Adminsignin)
	r.GET("/admin/viewuser", controllers.Adminviewuser)
	r.GET("/admin/adminpanel", controllers.Adminpanel)
	r.Run(os.Getenv("PORT"))

}
