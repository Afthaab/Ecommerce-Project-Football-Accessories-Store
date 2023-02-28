package controllers

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	var userdata models.User
	// DB := config.DBconnect()
	result = config.DB.Raw("SELECT firstname,lastname,email,phone,userid,created_at FROM users WHERE userid =?", id).Scan(&userdata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userdata models.User
	if c.Bind(&userdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Unable to Bind JSON data",
		})
		return
	}
	userdata.Userid = uint(id)
	// DB := config.DBconnect()
	result = config.DB.Model(&userdata).Updates(models.User{Firstname: userdata.Firstname, Lastname: userdata.Lastname, Phone: userdata.Phone}).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}
	// DB := config.DBconnect()
	result = config.DB.Raw("SELECT * from users WHERE userid = ?", id).Scan(&userdata).Error
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(userdata.Password), []byte(datas.Oldpassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	if datas.Password1 != datas.Password2 {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "Entered Password does not matches try again! ",
		})
		return
	}
	bytes, result1 := bcrypt.GenerateFromPassword([]byte(datas.Password1), 14)
	if result1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result1,
		})
		return
	}
	datas.Password1 = string(bytes)
	result2 := config.DB.Model(&models.User{}).Where("userid = ?", id).Update("password", datas.Password1).Error
	if result2 != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result2,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully Updated the Password",
	})

}

func ForgetPassword(c *gin.Context) {
	var userdata models.User
	if c.Bind(&userdata) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad Request",
		})
		return
	}
	result = config.DB.First(&userdata, "email LIKE ?", userdata.Email).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	otp := initializers.Otpgeneration(userdata.Email)
	result1 := config.DB.Model(&userdata).Where("email = ?", userdata.Email).Update("otp", otp)
	if result1.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Otp has been generated || Go to otp validate route",
	})
}

func ChangePassword(c *gin.Context) {
	type datas struct {
		Email    string
		Otp      string
		Password string
	}
	var signindata datas
	var userdata models.User
	var err error
	err = c.Bind(&signindata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	result = config.DB.First(&userdata, "email LIKE ? AND otp LIKE ?", signindata.Email, signindata.Otp).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	err = userdata.HashPassword(signindata.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Could not hash the password",
		})
		return
	}
	result1 := config.DB.Model(&userdata).Where("email = ?", userdata.Email).Update("password", userdata.Password)
	if result1.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
	})

}

func AdminProfilepage(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetString("adminid"))
	var admindata models.User
	// DB := config.DBconnect()
	result = config.DB.Raw("SELECT * FROM users WHERE userid = ?", id).Scan(&admindata).Error
	if result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Admin":       "Details",
		"First Name":  admindata.Firstname,
		"Last Name":   admindata.Lastname,
		"Email":       admindata.Email,
		"Phone":       admindata.Phone,
		"Is Admin":    admindata.IsAdmin,
		"Joined Date": admindata.CreatedAt,
	})
}
