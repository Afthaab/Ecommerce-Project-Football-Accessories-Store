package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Address struct {
	Addressid  uint   `JSON:"addressid" gorm:"primarykey;unique"`
	User       User   `gorm:"ForeignKey:uid"`
	Uid        uint   `JSON:"uid, omitempty"`
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
	ID         uint    `json:"id" gorm:"primaryKey"`
	Product    Product `gorm:"ForeignKey:Pid"`
	Pid        uint
	Quantity   uint
	Price      uint
	Totalprice uint
	User       User `gorm:"ForeignKey:Cartid"`
	Cartid     uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
	Baseimage   string `JSON:"baseimage"`
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
	ID            uint      `json:"id" gorm:"primaryKey"`
	Totalamount   uint      `JSON:"amount" gorm:"not null"`
	Paymentmethod string    `JSON:"paymentmethod" gorm:"not null"`
	Paymentstatus string    `JSON:"paymentstatus"`
	User          User      `gorm:"ForeignKey:Useridno"`
	Useridno      uint      `JSON:"useridno" gorm:"not null"`
	Coupon        Coupon    `gorm:"ForeignKey:Couponid"`
	Couponid      uuid.UUID `JSON:"couponid" gorm:"default:null"`
	RazorPay      RazorPay  `gorm:"ForeignKey:Razorpayid"`
	Razorpayid    string    `JSON:"razorpayid" gorm:"default:null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type RazorPay struct {
	UserID          uint   `JSON:"userid"`
	RazorPaymentId  string `JSON:"razorpaymentid" gorm:"primaryKey"`
	RazorPayOrderID string `JSON:"razorpayorderid"`
	Signature       string `JSON:"signature"`
	AmountPaid      string `JSON:"amountpaid"`
}
type Orders struct {
	Orderid     uuid.UUID `json:"orderid" gorm:"type:uuid;default:gen_random_uuid();not null;primaryKey"`
	User        User      `gorm:"ForeignKey:Useridno"`
	Useridno    uint      `json:"useridno"  gorm:"not null" `
	Totalamount uint      `json:"totalamount"  gorm:"not null" `
	Payment     Payment   `gorm:"ForeignKey:Paymentid"`
	Paymentid   uint      `json:"paymentid"`
	Address     Address   `gorm:"ForeignKey:Addid"`
	Addid       uint      `json:"addid"  `
	Orderstatus string    `json:"orderstatus" grom:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Orderditems struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Product    Product   `gorm:"ForeignKey:Pid"`
	Pid        uint      `json:"pid" gorm:"not null" `
	User       User      `gorm:"ForeignKey:Uid"`
	Uid        uint      `json:"uid"`
	Orders     Orders    `gorm:"ForeignKey:Oid"`
	Oid        uuid.UUID `json:"oid" gorm:"type:uuid;not null"`
	Quantity   uint      `json:"quantity" gorm:"not null"`
	Price      uint      `json:"price" gorm:"not null"`
	Totalprice uint      `json:"totalprice" gorm:"not null"`
}

type Coupon struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();not null;primaryKey"`
	Couponame      string    `json:"couponame" gorm:"not null"`
	Minamount      uint      `json:"minamount" gorm:"not null"`
	Discount       uint      `json:"discount" gorm:"not null"`
	Expirationdate time.Time `json:"expirationdate" gorm:"not null"`
	Isactive       bool      `json:"isactive" gorm:"not null;default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Wallet struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Balance   uint      `json:"balance" gorm:"not null"`
	User      User      `gorm:"ForeignKey:walletid"`
	Walletid  uint      `json:"walletid" gorm:"not null"`
	Orders    Orders    `gorm:"ForeignKey:Oid"`
	Oid       uuid.UUID `json:"oid" gorm:"not null"`
	CreatedAt time.Time `gorm:"default:NOW()"`
}

type Wishlist struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Product    Product   `gorm:"ForeignKey:Pid"`
	Pid        uint      `json:"pid" gorm:"not null" `
	User       User      `gorm:"ForeignKey:wishlistid"`
	Wishlistid uint      `json:"wishlistid" gorm:"not null"`
	CreatedAt  time.Time `gorm:"default:NOW()"`
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
