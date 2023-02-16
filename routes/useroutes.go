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

		//address routes
		user.POST("/address/add", middlewares.UserAuth, controllers.AddAddress)
		user.GET("/address/view", middlewares.UserAuth, controllers.ShowAddress)
		user.PUT("/address/edit", middlewares.UserAuth, controllers.EditAddress)

		//cart routes
		user.POST("/cart/add", middlewares.UserAuth, controllers.AddToCart)
		user.GET("/cart/view", middlewares.UserAuth, controllers.ViewCart)
		user.PUT("/cart/view/edit", middlewares.UserAuth, controllers.EditCart)
		user.DELETE("/cart/view/delete", middlewares.UserAuth, controllers.DeleteCartItem)

		//checkout routes
		user.GET("/checkoutpage", middlewares.UserAuth, controllers.CheckOut)

		//coupon routes
		user.POST("/coupon/redeem", middlewares.UserAuth, controllers.RedeemCoupon)

		//payments route
		user.POST("/payment/cod", middlewares.UserAuth, controllers.CashOnDevlivery)
		user.GET("/payment/razorpay", middlewares.UserAuth, controllers.Razorpay)
		user.GET("/payment/success", middlewares.UserAuth, controllers.RazorpaySuccess)
		user.GET("/success", middlewares.UserAuth, controllers.Success)

		//order routes
		user.POST("/order/placeorder", middlewares.UserAuth, controllers.PlaceOrder)
		user.GET("/order/view/search", middlewares.UserAuth, controllers.ViewOrder)
		user.PUT("/order/cancel", middlewares.UserAuth, controllers.CancelOrder)
		user.PUT("/order/return", middlewares.UserAuth, controllers.ReturnOrder)
		user.GET("/order/view/invoice", middlewares.UserAuth, controllers.GenerateInvoice)
		user.GET("/order/view/invoice/download", middlewares.UserAuth, controllers.InvoiceDownload)

		//Wishlist routes
		user.POST("/product/view/wishlist/add", middlewares.UserAuth, controllers.AddToWishlist)
		user.GET("/product/view/wishlist/view", middlewares.UserAuth, controllers.ViewWishlist)
		user.DELETE("/product/view/wishlist/delete", middlewares.UserAuth, controllers.RemoveWishlist)

		//wallet
		user.GET("/wallet/view", middlewares.UserAuth, controllers.ShowWallet)
	}

}
