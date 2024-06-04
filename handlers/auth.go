package handlers

import (
	"encoding/json"
	"finalProjStart/entity"
	"finalProjStart/error"
	"finalProjStart/jsonlog"
	"finalProjStart/repository"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	repo   repository.UserRepository // Declare a global variable for UserRepository
	logger *jsonlog.Logger
)

const SecretKey = "secret"

func InitUserRepository(userRepo repository.UserRepository) {
	repo = userRepo // Assign the provided userRepo to the global repo variable
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.PrintError(err, map[string]string{"context": "decoding user"})
		errors.HandleBadRequest(w, "Invalid request payload")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.PrintError(err, map[string]string{"context": "generating hashed password"})
		errors.HandleInternalServerError(w, err)
		return
	}
	user.Password = string(hashedPassword)
	user.Role = "user" // Assigning the role "user" by default
	user.RegisteredAt = time.Now()

	err = repo.SaveUser(&user)
	if err != nil {
		logger.PrintError(err, map[string]string{"context": "saving user"})
		errors.HandleInternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		logger.PrintError(err, map[string]string{"context": "decoding credentials"})
		errors.HandleBadRequest(w, "Invalid request payload")
		return
	}

	user, err := repo.FindUserByEmail(credentials.Email)
	if err != nil {
		logger.PrintError(err, map[string]string{"context": "finding user by email"})
		errors.HandleUnauthorized(w)
		return
	}

	if user == nil {
		logger.PrintError(err, map[string]string{"context": "user not found"})
		errors.HandleUnauthorized(w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		logger.PrintError(err, map[string]string{"context": "comparing passwords"})
		errors.HandleUnauthorized(w)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logger.PrintError(err, map[string]string{"context": "generating JWT token"})
		errors.HandleInternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Authorization", "")

	// Clear the token cookie from the client-side
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1, // MaxAge<0 means delete cookie now
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true, // Set to true if using HTTPS
		Path:     "/",
	})

	logger.PrintInfo("User logged out", nil)

	// Send a JSON response indicating successful logout
	response := map[string]string{"message": "Logged out successfully"}
	json.NewEncoder(w).Encode(response)
}
