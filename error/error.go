package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIError represents a structured error response.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func SendErrorResponse(w http.ResponseWriter, apiErr *APIError) {
	w.WriteHeader(apiErr.StatusCode)
	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		log.Printf("Failed to send error response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// HandleNotFound handles 404 Not Found errors.
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	SendErrorResponse(w, NewAPIError(http.StatusNotFound, "Resource not found"))
}

// HandleMethodNotAllowed handles 405 Method Not Allowed errors.
func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	SendErrorResponse(w, NewAPIError(http.StatusMethodNotAllowed, "Method not allowed"))
}

// HandleInternalServerError handles 500 Internal Server Error.
func HandleInternalServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v", err)
	SendErrorResponse(w, NewAPIError(http.StatusInternalServerError, "Internal server error"))
}

// HandleBadRequest handles 400 Bad Request errors.
func HandleBadRequest(w http.ResponseWriter, message string) {
	SendErrorResponse(w, NewAPIError(http.StatusBadRequest, message))
}

// HandleUnauthorized handles 401 Unauthorized errors.
func HandleUnauthorized(w http.ResponseWriter) {
	SendErrorResponse(w, NewAPIError(http.StatusUnauthorized, "Unauthorized"))
}

// HandleForbidden handles 403 Forbidden errors.
func HandleForbidden(w http.ResponseWriter) {
	SendErrorResponse(w, NewAPIError(http.StatusForbidden, "Forbidden"))
}
