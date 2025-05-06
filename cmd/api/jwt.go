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

const (
	UserKey   contextKey = "UserID"
	SellerKey contextKey = "SellerID"
	AdminKey  contextKey = "AdminID"
)

// JWTCreate generates a JWT token with the provided secret, user ID, and type.
func JWTCreate(secret []byte, userID int) (string, error) {
	expiration := time.Duration(env.GetInt("JWTExpirationInSeconds", 3600)) * time.Second

	claims := jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateToken parses and validates the provided JWT token string.
func validateToken(token string) (*jwt.Token, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(env.GetString("JWTSECRET", "secret")), nil
	})
}

// GetTokenFromRequest extracts the JWT token from the request's Authorization header.
func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if len(tokenAuth) > 7 && tokenAuth[:7] == "Bearer " {
		return tokenAuth[7:]
	}
	return ""
}
