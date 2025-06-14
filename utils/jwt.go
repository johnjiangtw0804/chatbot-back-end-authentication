package utils

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(secretKey string, userID uint, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt: jwt.NewNumericDate(time.Now())},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// HMAC signing method2	HS256,HS384,HS512	[]byte	[]byte
	return token.SignedString([]byte(secretKey))
}
