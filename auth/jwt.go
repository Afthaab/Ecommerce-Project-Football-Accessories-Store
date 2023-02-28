package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TokenGeneration(id string) (map[string]string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": id,
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		panic(err)
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["userid"] = id
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"access_token":  tokenString,
		"expiry_time":   "1 hr",
		"refresh_token": rt,
	}, nil
}
