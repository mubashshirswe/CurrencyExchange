package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mubashshir3767/currencyExchange/internal/store"
)

func (app *application) JWTUserMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := GetTokenFromRequest(r)
			token, err := ValidateToken(tokenString)
			if err != nil {
				log.Printf("failed to validate token: %v", err)
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			if !token.Valid {
				log.Println("invalid token")
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid token"))
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			expirationTime := int64(claims["expiredAt"].(float64))
			if expirationTime < time.Now().Unix() {
				log.Println("token has expired")
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("token has expired"))
				return
			}

			userString := claims["userID"].(string)
			userID, _ := strconv.Atoi(userString)

			user, err := app.GetUser(int64(userID))
			if err != nil {
				log.Printf("failed to get user by id %v", err)
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (app *application) GetUser(id int64) (*store.User, error) {
	user, err := app.cacheStore.Users.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = app.store.Users.GetById(context.Background(), &id)
		if err != nil {
			return nil, err
		}

		if err := app.cacheStore.Users.Set(context.Background(), user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
