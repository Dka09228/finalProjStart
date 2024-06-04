package middleware

import (
	"finalProjStart/error"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

var SecretKey = "secret"

func AuthMiddleware(next http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.HandleUnauthorized(w)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil {
			errors.HandleUnauthorized(w)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] != role {
				errors.HandleForbidden(w)
				return
			}
			next(w, r)
		} else {
			errors.HandleUnauthorized(w)
			return
		}
	}
}
