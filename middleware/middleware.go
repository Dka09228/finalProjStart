package middleware

import (
	errors "finalProjStart/error"
	"finalProjStart/jsonlog"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

var SecretKey = "secret"
var logger *jsonlog.Logger

func init() {
	logger = jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)
}

func AuthMiddleware(next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.HandleUnauthorized(w)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			errors.HandleUnauthorized(w)
			return
		}

		tokenString := headerParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			errors.HandleUnauthorized(w)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if !containsRole(roles, claims["role"].(string)) {
				errors.HandleForbidden(w)
				return
			}
		} else {
			errors.HandleUnauthorized(w)
			return
		}

		next(w, r)
	}
}

func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logger.PrintPanic(err, stack)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
