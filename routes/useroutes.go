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
		user.POST("/forgetpassword", controllers.ForgetPassword)
		user.POST("/forgetpassword/validateotp", controllers.ValidateOtp)
		user.POST("/forgetpassword/changepassword", controllers.ChangePassword)
		user.GET("/signout", middlewares.UserAuth, controllers.UserSignout)

		//profile page routes
		user.GET("/profilepage", middlewares.UserAuth, controllers.GetUserProfile)
		user.PUT("/profilepage/editprofile", middlewares.UserAuth, controllers.EditUserProfile)
		user.PUT("/profilepage/editprofile/changepassword", middlewares.UserAuth, controllers.ChangePasswordInProfile)
		user.POST("/profilepage/addaddress", middlewares.UserAuth, controllers.AddAddress)
		user.GET("/profilepage/showaddress", middlewares.UserAuth, controllers.ShowAddress)
		user.PUT("/profilepage/editaddress", middlewares.UserAuth, controllers.EditAddress)

		//adding product to the cart
		user.POST("/addtocart", middlewares.UserAuth, controllers.AddToCart)
		user.GET("/viewcart", middlewares.UserAuth, controllers.ViewCart)

		user.GET("/checkoutpage", middlewares.UserAuth, controllers.CheckOut)

	}

}
