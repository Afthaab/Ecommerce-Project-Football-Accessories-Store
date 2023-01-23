package routes

import (
	"github.com/afthab/e_commerce/controllers"
	"github.com/afthab/e_commerce/middlewares"
	"github.com/gin-gonic/gin"
)

func Userroutes(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/signup", controllers.Usersignup)
		user.POST("/signup/otpvalidate", controllers.Otpvalidate)
		user.POST("/signin", controllers.Usersignin)

		//routes with middlewares
		user.GET("/profilepage", middlewares.UserAuth, controllers.GetUserProfile)
		user.POST("/address", middlewares.UserAuth, controllers.AddAddress)
		user.GET("/profilepage/showaddress", middlewares.UserAuth, controllers.ShowAddress)
		
	}

}
