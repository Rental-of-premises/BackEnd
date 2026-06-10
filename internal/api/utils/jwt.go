package utils

import (
	"os"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    ID int64   `json:"id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateJWT(id int64, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET не установлен в переменных окружения")
	}

	claims := &Claims{
        ID: id,
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt: jwt.NewNumericDate(time.Now()),
            Issuer: "rent-api",
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}