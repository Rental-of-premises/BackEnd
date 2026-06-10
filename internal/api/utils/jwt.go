package utils

import (
	"time"
	"errors"
	"strings"
	"net/http"

	config "rent/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    ID int64   `json:"id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateJWT(id int64, email string) (string, error) {
	cnf := config.GetSingletonConfig()
	secret := cnf.JWTSecret
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

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
            return nil, errors.New("unexpected signing method")
        }
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}

func ExtractToken(r *http.Request) string {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return ""
    }
    
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return ""
    }
    
    return parts[1]
}