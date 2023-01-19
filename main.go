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
	r.POST("/signup", controllers.Usersignup)
	r.POST("/signin", controllers.Usersignin)
	r.Run(os.Getenv("PORT"))

}
