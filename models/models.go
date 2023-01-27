package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Userid      uint   `JSON:"userid" gorm:"unique;primaryKey"`
	Firstname   string `JSON:"firstname" gorm:"not_null" validate:"required,min=2,max=50"`
	Lastname    string `JSON:"lastname" gorm:"not null" validate:"required,min=2,max=50"`
	Email       string `JSON:"email" gorm:"not null;unique;" validate:"email,required" `
	Phone       string `JSON:"phone" gorm:"not null;unique" validate:"required,len=10"`
	Password    string `JSON:"password" gorm:"not null" validate:"required"`
	Otpverified bool   `JSON:"otpverified" gorm:"default:false"`
	Isblocked   bool   `JSON:"isblocked" gorm:"default:false"`
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
	Userid    uint   `JSON:"userid" gorm:"foreignKey:UserRefer"`
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
type Product struct {
	Productid   uint   `JSON:"productid" gorm:"primarykey;unique"`
	Productname string `JSON:"productname" gorm:"not null"`
	Description string `JSON:"description" gorm:"not null"`
	Stock       uint   `JSON:"stock" gorm:"not null"`
	Price       uint   `JSON:"price" gorm:"not null"`
	Teamid      uint   `JSON:"teamid"`
	Brandid     uint   `JSON:"brandid"`
	Sizeid      uint   `JSON:"sizeid"`
}

type Size struct {
	Product  Product `gorm:"foreignkey:Sizeid"`
	Sizeid   uint    `JSON:"sizeid" gorm:"primarykey"`
	Sizetype string  `JSON:"sizetype" gorm:"not null"`
}

type Brand struct {
	Product   Product `gorm:"foreignkey:Brandid"`
	Brandid   uint    `JSON:"brandid" gorm:"primarykey"`
	Brandname string  `JSON:"brandname" gorm:"not null"`
}

type Team struct {
	Product  Product `gorm:"foreignkey:Teamid"`
	Teamid   uint    `JSON:"teamid" gorm:"primarykey"`
	Teamname string  `JSON:"teamname" gorm:"not null"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		panic(err)
	}
	return nil
}
