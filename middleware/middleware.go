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
			logger.PrintInfo("Authorization header missing", nil)
			errors.HandleUnauthorized(w)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			logger.PrintInfo("Invalid Authorization header format", nil)
			errors.HandleUnauthorized(w)
			return
		}

		tokenString := headerParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			logger.PrintInfo("Invalid token", map[string]string{"error": err.Error()})
			errors.HandleUnauthorized(w)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role, ok := claims["role"].(string)
			if !ok {
				logger.PrintInfo("Role claim missing or invalid", nil)
				errors.HandleUnauthorized(w)
				return
			}

			role = strings.TrimSpace(role) // Trim whitespace from the role
			logger.PrintInfo("Token role", map[string]string{"role": role})

			if !containsRole(roles, role) {
				logger.PrintInfo("Role not permitted", map[string]string{"role": role})
				errors.HandleForbidden(w)
				return
			}

			logger.PrintInfo("Role permitted", map[string]string{"role": role})
		} else {
			logger.PrintInfo("Token claims not valid", nil)
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
