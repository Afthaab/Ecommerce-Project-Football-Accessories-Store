package routes

import (
	"github.com/afthab/e_commerce/controllers"
	"github.com/gin-gonic/gin"
)

func Userroutes(r *gin.Engine) {
	r.POST("/user/signup", controllers.Usersignup)
	r.POST("/user/signup/otpvalidate", controllers.Otpvalidate)
	r.POST("/user/signin", controllers.Usersignin)
}
