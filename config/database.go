package config

import (
	"fmt"
	"os"

	"github.com/afthab/e_commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBconnect() *gorm.DB {
	dns := os.Getenv("DB_URL")
	DB, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		fmt.Println("Error Connecting to the base")
	}
	DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Address{},
		&models.Brand{},
		&models.Team{},
		&models.Size{},
		&models.Product{},
		&models.Image{},
		&models.Cart{},
	)
	return DB
}
