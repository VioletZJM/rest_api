package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

type AppClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func VerifyToken(token string) (int64, error) {
	claims := &AppClaims{}

	parseToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if !parseToken.Valid {
		return 0, errors.New("Invalid token")
	}

	return claims.UserId, nil
}

func LoadJWTSecret() error {
	raw, ok := os.LookupEnv("JWT_SECRET")

	if !ok || raw == "" {
		return errors.New("JWT_SECRET is not set")
	}

	jwtSecret = []byte(raw)

	return nil
}
