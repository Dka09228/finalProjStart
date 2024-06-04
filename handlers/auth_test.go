package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"finalProjStart/entity"
	"finalProjStart/handlers"
	"finalProjStart/repository"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mockRepo repository.UserRepository
	testDB   *mongo.Database
)

func setupMongoDB() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database("test_db")
	return db, nil
}

func setupTest() {
	var err error
	testDB, err = setupMongoDB()
	if err != nil {
		log.Fatalf("Failed to set up MongoDB: %v", err)
	}

	mockRepo = repository.NewMongoDBUserRepository(testDB)
	handlers.InitUserRepository(mockRepo)
}

func teardownTest() {
	if testDB != nil {
		err := testDB.Drop(context.Background())
		if err != nil {
			log.Printf("Error dropping test database: %v", err)
		}
	}
}

func TestRegisterUserHandlerIntegration(t *testing.T) {
	setupTest()
	defer teardownTest()

	user := entity.User{
		Email:    "test@example.com",
		Password: "password",
	}

	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.RegisterUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser entity.User
	err = json.NewDecoder(rr.Body).Decode(&responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, responseUser.Email)
	assert.NotEmpty(t, responseUser.ID.Hex()) // Ensure ID is assigned and converted to hex
}

func TestLoginUserHandlerIntegration(t *testing.T) {
	setupTest()
	defer teardownTest()

	// Register a user first
	registerUser := entity.User{
		Email:    "test@example.com",
		Password: "password",
	}
	registerUserJSON, err := json.Marshal(registerUser)
	assert.NoError(t, err)

	registerReq, err := http.NewRequest("POST", "/register", bytes.NewBuffer(registerUserJSON))
	assert.NoError(t, err)

	registerRR := httptest.NewRecorder()
	registerHandler := http.HandlerFunc(handlers.RegisterUser)
	registerHandler.ServeHTTP(registerRR, registerReq)
	assert.Equal(t, http.StatusOK, registerRR.Code)

	// Now test login with registered credentials
	loginCredentials := map[string]string{
		"email":    "test@example.com",
		"password": "password",
	}
	loginCredentialsJSON, err := json.Marshal(loginCredentials)
	assert.NoError(t, err)

	loginReq, err := http.NewRequest("POST", "/login", bytes.NewBuffer(loginCredentialsJSON))
	assert.NoError(t, err)

	loginRR := httptest.NewRecorder()
	loginHandler := http.HandlerFunc(handlers.LoginUser)
	loginHandler.ServeHTTP(loginRR, loginReq)

	assert.Equal(t, http.StatusOK, loginRR.Code)

	var response map[string]string
	err = json.NewDecoder(loginRR.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}
	assert.NotEmpty(t, response["token"])
}

func TestRegisterUserSuccessIntegration(t *testing.T) {
	setupTest()
	defer teardownTest()

	// Create valid user data
	user := entity.User{
		Email:    "test@example.com",
		Password: "password",
	}

	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)

	// Create HTTP request with valid user data
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)

	// Perform the request
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.RegisterUser)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 for successful registration")

	// Check the response body for the registered user
	var responseUser entity.User
	err = json.NewDecoder(rr.Body).Decode(&responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, responseUser.Email)
	assert.NotEmpty(t, responseUser.ID.Hex()) // Ensure ID is assigned and converted to hex
}
