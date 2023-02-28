package config

import (
	"fmt"
	"os"

	"github.com/afthab/e_commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBconnect() {
	var err error
	dns := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		fmt.Println("Error Connecting to the base")
	}
	DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Brand{},
		&models.Team{},
		&models.Size{},
		&models.Product{},
		&models.Image{},
		&models.Cart{},
		&models.Payment{},
		&models.Orders{},
		&models.Orderditems{},
		&models.Coupon{},
		&models.RazorPay{},
		&models.Wallet{},
		&models.Wishlist{},
	)
}
