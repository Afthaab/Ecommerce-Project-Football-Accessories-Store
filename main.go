package main

import (
	"os"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/initializers"
	"github.com/afthab/e_commerce/routes"
	"github.com/gin-gonic/gin"
)

var R = gin.Default()

func init() {
	initializers.LoadEnv()
	config.DBconnect()
	R.LoadHTMLGlob("templates/*.html")
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	routes.Userroutes(R)
	routes.Adminroutes(R)
	routes.OpenRoutes(R)
	R.Run(os.Getenv("PORT"))
}
