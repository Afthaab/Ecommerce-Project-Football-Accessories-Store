package models

import (
	"time"
)

type User struct {
	Userid      uint   `JSON:"userid" gorm:"unique;primaryKey"`
	Firstname   string `JSON:"firstname" gorm:"not_null" validate:"required,min=2,max=50"`
	Lastname    string `JSON:"lastname" gorm:"not null" validate:"required,min=2,max=50"`
	Email       string `JSON:"email" gorm:"not null;unique;" validate:"email,required" `
	Phone       string `JSON:"phone" gorm:"not null;unique" validate:"required,len=10"`
	Password    string `JSON:"password" gorm:"not null" validate:"required"`
	Otpverified bool   `JSON:"otpverified" gorm:"default:false"`
	Isblocked   bool   `JSON:"isblocked" gorm:"default:true"`
	Otp         string `JSON:"otp"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Admin struct {
	Adminid   uint
	Firstname string `JSON:"firstname"`
	Lastname  string `JSON:"lastname"`
	Email     string `JSON:"email"`
	Phone     string `JSON:"phone"`
	Password  string `JSON:"password"`
}
type Address struct {
	Addressid uint   `JSON:"addressid" gorm:"primarykey;unique"`
	Userid    uint   `JSON:"userid" gorm:"foriegnkey"`
	Name      string `JSON:"name" gorm:"not null"`
	Phoneno   string `JSON:"phoneno" gorm:"not null"`
	Houseno   string `JSON:"houseno" gorm:"not null"`
	Area      string `JSON:"area" gorm:"not null"`
	Landmark  string `JSON:"landmark" gorm:"not null"`
	City      string `JSON:"city" gorm:"not null"`
	Pincode   string `JSON:"pincode" gorm:"not null"`
	District  string `JSON:"district" gorm:"not null"`
	State     string `JSON:"state" gorm:"not null"`
	Country   string `JSON:"country" gorm:"not null"`
}
