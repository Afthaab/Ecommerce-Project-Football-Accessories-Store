package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TokenGeneration(id string) string {
	//creating a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenstring, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		panic(err)
	}
	return tokenstring
}
