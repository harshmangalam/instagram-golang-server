package utils

import (
	"instagram/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWTToken(userId int) (string, error) {

	type TokenClaim struct {
		UserId int `json:"userId"`
		jwt.StandardClaims
	}

	claims := TokenClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			Issuer:    "token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString([]byte(config.Get("JWT_SECRET")))

	return str, err

}
