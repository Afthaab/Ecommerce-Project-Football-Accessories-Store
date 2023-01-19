package config

import (
	"log"
	"os"

	"github.com/afthab/e_commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBconnect() *gorm.DB{
	dns:=os.Getenv("DB_URL")
	DB, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil{
		log.Fatal(err)
	}
	DB.AutoMigrate(&models.User{})
	return DB
}