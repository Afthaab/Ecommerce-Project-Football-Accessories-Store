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
		user.POST("/forgetpassword/changepassword", controllers.ChangePassword)
		user.GET("/signout", middlewares.UserAuth, controllers.UserSignout)

		//profile page routes
		profilepage := user.Group("/profilepage")
		{
			profilepage.GET("/", middlewares.UserAuth, controllers.GetUserProfile)
			profilepage.PUT("/editprofile", middlewares.UserAuth, controllers.EditUserProfile)
			profilepage.PUT("/editprofile/changepassword", middlewares.UserAuth, controllers.ChangePasswordInProfile)
		}

		//address routes
		address := user.Group("/address")
		{
			address.POST("/add", middlewares.UserAuth, controllers.AddAddress)
			address.GET("/view", middlewares.UserAuth, controllers.ShowAddress)
			address.PUT("/edit", middlewares.UserAuth, controllers.EditAddress)
		}

		//cart routes
		cart := user.Group("/cart")
		{
			cart.POST("/add", middlewares.UserAuth, controllers.AddToCart)
			cart.GET("/view", middlewares.UserAuth, controllers.ViewCart)
			cart.PUT("/view/edit", middlewares.UserAuth, controllers.EditCart)
			cart.DELETE("/view/delete", middlewares.UserAuth, controllers.DeleteCartItem)

		}

		//checkout routes
		user.GET("/checkoutpage", middlewares.UserAuth, controllers.CheckOut)

		//coupon routes
		user.POST("/coupon/redeem", middlewares.UserAuth, controllers.RedeemCoupon)

		//payments route
		payment := user.Group("/payment")
		{
			payment.POST("/cod", middlewares.UserAuth, controllers.CashOnDevlivery)
			payment.GET("/razorpay", middlewares.UserAuth, controllers.Razorpay)
			payment.GET("/success", middlewares.UserAuth, controllers.RazorpaySuccess)
			payment.GET("/successpage", middlewares.UserAuth, controllers.Success)
			payment.POST("/wallet", middlewares.UserAuth, controllers.WalletPayment)

		}

		//order routes
		order := user.Group("/order")
		{
			order.POST("/placeorder", middlewares.UserAuth, controllers.PlaceOrder)
			order.GET("/view/search", middlewares.UserAuth, controllers.ViewOrder)
			order.PUT("/cancel", middlewares.UserAuth, controllers.CancelOrder)
			order.PUT("/return", middlewares.UserAuth, controllers.ReturnOrder)
			order.GET("/view/invoice", middlewares.UserAuth, controllers.GenerateInvoice)
			order.GET("/view/invoice/download", middlewares.UserAuth, controllers.InvoiceDownload)
		}

		//Wishlist routes
		wishlist := user.Group("/product/view/wishlist")
		{
			wishlist.POST("/add", middlewares.UserAuth, controllers.AddToWishlist)
			wishlist.GET("/view", middlewares.UserAuth, controllers.ViewWishlist)
			wishlist.DELETE("/delete", middlewares.UserAuth, controllers.RemoveWishlist)
		}

		//wallet
		user.GET("/wallet/view", middlewares.UserAuth, controllers.ShowWallet)
	}

}
