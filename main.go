package main

import (
	"os"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/initializers"
	"github.com/afthab/e_commerce/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	config.DBconnect()
}

var R = gin.Default()

func main() {
	gin.SetMode(gin.ReleaseMode)

	routes.Userroutes(R)
	routes.Adminroutes(R)

	R.Run(os.Getenv("PORT"))

}
