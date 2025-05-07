package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mubashshir3767/currencyExchange/internal/env"
)

type contextkey string

const UserKey contextkey = "UserID"
const SellerKey contextkey = "SellerID"
const AdminKey contextkey = "AdminID"

func JWTCreate(secret []byte, userID int, idType string) (string, error) {
	expiration := time.Second * time.Duration(env.GetInt("JWTExpirationInSeconds", 3600000))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		idType:      strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(token string) (*jwt.Token, error) {
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
		return tokenAuth[7:]
	}
	return ""
}
