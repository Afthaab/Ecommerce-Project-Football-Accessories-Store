package models

import (
	"time"
)

type User struct {
	Userid      uint   `JSON:"userid" gorm:"unique;primaryKey"`
	Firstname   string `JSON:"firstname" gorm:"not_null" validate:"required,min=2,max=50"`
	Lastname    string `JSON:"lastname" gorm:"not null" validate:"required,min=2,max=50"`
	Email       string `JSON:"email" gorm:"not null;unique" validate:"email,required" `
	Phone       string `JSON:"phone" gorm:"not null;unique" validate:"required,len=10"`
	Password    string `JSON:"password" gorm:"not null" validate:"required"`
	Otpverified bool   `JSON:"otpverified" gorm:"default:false"`
	Isblocked   bool   `JSON:"Isblocked" gorm:"default:false"`
	Otp         string   `JSON:"otp"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
