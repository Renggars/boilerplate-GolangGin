package utils

import (
	"os"
	"restApi-GoGin/src/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))
var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

type JWTAccessClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTRefreshClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *models.User) (string, error) {
	claims := JWTAccessClaims{
		user.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(accessSecret)

	return ss, err
}

func GenerateRefreshToken(user *models.User) (string, error) {
	claims := JWTRefreshClaims{
		user.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(refreshSecret)

	return ss, err
}

func VerifyRefreshToken(tokenStr string) (*JWTRefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTRefreshClaims{}, func(t *jwt.Token) (interface{}, error) { return refreshSecret, nil })

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTRefreshClaims)

	if !ok {
		return nil, err
	}

	return claims, nil

}

func VerifyAccessToken(tokenStr string) (*JWTAccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) { return accessSecret, nil })

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTAccessClaims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
