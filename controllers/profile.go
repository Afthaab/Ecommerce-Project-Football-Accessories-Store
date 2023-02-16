package controllers

import (
	"strconv"

	"github.com/afthab/e_commerce/config"
	"github.com/afthab/e_commerce/initializers"
	"github.com/afthab/e_commerce/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	var userdata models.User
	DB := config.DBconnect()
	result := DB.Raw("SELECT firstname,lastname,email,phone,userid,created_at FROM users WHERE userid =?", id).Scan(&userdata)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"User ID":     userdata.Userid,
		"First Name":  userdata.Firstname,
		"Last Name":   userdata.Lastname,
		"Email":       userdata.Email,
		"Phone":       userdata.Phone,
		"Date Joined": userdata.CreatedAt,
	})

}

func EditUserProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userdata models.User
	if c.Bind(&userdata) != nil {
		c.JSON(400, gin.H{
			"Error": "Unable to Bind JSON data",
		})
		return
	}
	userdata.Userid = uint(id)
	DB := config.DBconnect()
	result := DB.Model(&userdata).Updates(models.User{Firstname: userdata.Firstname, Lastname: userdata.Lastname, Phone: userdata.Phone})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Profile Updated Successfully",
		"User ID": userdata.Userid,
	})

}

func ChangePasswordInProfile(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("userid"))
	type passwordata struct {
		Oldpassword string
		Password1   string
		Password2   string
	}
	var datas passwordata
	var userdata models.User
	if c.Bind(&datas) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Raw("SELECT * from users WHERE userid = ?", id).Scan(&userdata).Error
	if result != nil {
		c.JSON(400, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(userdata.Password), []byte(datas.Oldpassword))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	if datas.Password1 != datas.Password2 {
		c.JSON(404, gin.H{
			"Error": "Entered Password does not matches try again! ",
		})
		return
	}
	bytes, result1 := bcrypt.GenerateFromPassword([]byte(datas.Password1), 14)
	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": result1,
		})
		return
	}
	datas.Password1 = string(bytes)
	result2 := DB.Model(&models.User{}).Where("userid = ?", id).Update("password", datas.Password1).Error
	if result2 != nil {
		c.JSON(404, gin.H{
			"Error": result2,
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully Updated the Password",
	})

}

func ForgetPassword(c *gin.Context) {
	type email struct {
		Email string
	}
	var emaildata email
	var userdata models.User
	if c.Bind(&emaildata) != nil {
		c.JSON(400, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.First(&userdata, "email LIKE ?", emaildata.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	otp := initializers.Otpgeneration(userdata.Email)
	result1 := DB.Model(&userdata).Where("email = ?", emaildata.Email).Update("otp", otp)
	if result1.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"Error": "User not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Otp has been generated || Go to otp validate route",
	})
}

func ChangePassword(c *gin.Context) {
	var userdata models.User
	if c.Bind(&userdata) != nil {
		c.JSON(400, gin.H{
			"Errro": "Error in binding the JSON data",
		})
		return
	}
	err := userdata.HashPassword(userdata.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Could not hash the password",
		})
		return
	}
	DB := config.DBconnect()
	result1 := DB.Model(&userdata).Where("email = ?", userdata.Email).Update("password", userdata.Password)
	if result1.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"Error": "User not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Success",
	})

}

func ValidateOtp(c *gin.Context) {
	type Userotp struct {
		Otp   string
		Email string
	}
	var otpdata Userotp
	DB := config.DBconnect()
	result := DB.Raw("SELECT * FROM users WHERE email LIKE ? AND opt LIKE ?", otpdata.Email, otpdata.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "OTP Validated",
	})

}

func AdminProfilepage(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("adminid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var admindata models.Admin
	DB := config.DBconnect()
	result := DB.Raw("SELECT firstname,lastname,email,phone FROM admins WHERE adminid = ?", id).Scan(&admindata)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Admin Details": admindata,
	})
}
