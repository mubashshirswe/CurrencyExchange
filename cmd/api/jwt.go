package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mubashshir3767/currencyExchange/internal/env"
)

type contextKey string

const UserKey contextKey = "UserID"
const EmployeeKey contextKey = "EmployeeID"

func JWTCreate(secret []byte, id int64, types string) (string, error) {
	expiration := time.Second * time.Duration(env.GetInt("JWTExpirationInSeconds", 3600000))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		types:       strconv.Itoa(int(id)),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(token)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(token string) (*jwt.Token, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(env.GetString("JWTSECRET", "secret")), nil
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		return tokenAuth[:7]
	}

	return ""
}
