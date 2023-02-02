package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Address struct {
	Addressid  uint   `JSON:"addressid" gorm:"primarykey;unique"`
	User       User   `gorm:"ForeignKey:uid"`
	Uid        uint   `JSON:"uid"`
	Name       string `JSON:"name" gorm:"not null"`
	Phoneno    string `JSON:"phoneno" gorm:"not null"`
	Houseno    string `JSON:"houseno" gorm:"not null"`
	Area       string `JSON:"area" gorm:"not null"`
	Landmark   string `JSON:"landmark" gorm:"not null"`
	City       string `JSON:"city" gorm:"not null"`
	Pincode    string `JSON:"pincode" gorm:"not null"`
	District   string `JSON:"district" gorm:"not null"`
	State      string `JSON:"state" gorm:"not null"`
	Country    string `JSON:"country" gorm:"not null"`
	Defaultadd bool   `JSON:"defaultadd" gorm:"default:false"`
}

type Cart struct {
	gorm.Model
	Product    Product `gorm:"ForeignKey:Pid"`
	Pid        uint
	Quantity   uint
	Price      uint
	Totalprice uint
	User       User `gorm:"ForeignKey:Cartid"`
	Cartid     uint
}

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

type Image struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	Product Product `gorm:"ForeignKey:Pid"`
	Pid     uint    `json:"pid"`
	Image   string  `JSON:"Image" gorm:"not null"`
}

type Product struct {
	Productid   uint   `JSON:"productid" gorm:"primarykey;unique"`
	Productname string `JSON:"productname" gorm:"not null"`
	Description string `JSON:"description" gorm:"not null"`
	Stock       uint   `JSON:"stock" gorm:"not null"`
	Price       uint   `JSON:"price" gorm:"not null"`
	Team        Team   `gorm:"ForeignKey:Teamid"`
	Teamid      uint   `JSON:"teamid"`
	Brand       Brand  `gorm:"ForeignKey:Brandid"`
	Brandid     uint   `JSON:"brandid"`
	Size        Size   `gorm:"ForeignKey:Sizeid"`
	Sizeid      uint   `JSON:"sizeid"`
}

type Size struct {
	ID       uint   `json:"id" gorm:"primaryKey"  `
	Sizetype string `JSON:"sizetype" gorm:"not null"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primaryKey"  `
	Brandname string `JSON:"brandname" gorm:"not null"`
}

type Team struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Teamname string `JSON:"teamname" gorm:"not null"`
}

type Payment struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Totalamount   uint   `JSON:"amount" gorm:"not null"`
	Paymentmethod string `JSON:"paymentmethod" gorm:"not null"`
	Paymentstatus bool   `JSON:"paymentstatus" gorm:"defualt:flase"`
	User          User   `gorm:"ForeignKey:Uid"`
	Uid           uint   `JSON:"uid" gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type Orders struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	User        User    `gorm:"ForeignKey:Uid"`
	Uid         uint    `json:"uid"  gorm:"not null" `
	Totalamount uint    `json:"totalamount"  gorm:"not null" `
	Payment     Payment `gorm:"ForeignKey:pid"`
	Pid         uint    `json:"pid"`
	Orderstatus string  `json:"orderstatus"   `
	Address     Address `gorm:"ForeignKey:Addid"`
	Addid       uint    `json:"addid"  `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// type Orderditems struct {
// 	gorm.Model
// 	Uid         uint `json:"uid"  gorm:"not null" `
// 	Pid         uint `json:"pid" gorm:"not null" `
// 	Orders      Orders `gorm:"ForeignKey:Orderid"`
// 	Orderid     string `json:"orderid" gorm:"not null" `
// 	Orderstatus string `json:"Orderstatus" gorm:"not null" `
// }

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
