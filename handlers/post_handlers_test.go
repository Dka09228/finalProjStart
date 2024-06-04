package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"finalProjStart/entity"
	"finalProjStart/handlers"
	"finalProjStart/repository"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDBPosts *mongo.Database // MongoDB test database instance for posts
)

func TestMain(m *testing.M) {
	setupDBPosts()
	code := m.Run()
	teardownDBPosts()
	os.Exit(code)
}

func setupDBPosts() {
	// Connect to MongoDB test instance
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	testDBPosts = client.Database("test_db_posts") // Use a different database name
}

func teardownDBPosts() {
	// Drop the test database after tests
	err := testDBPosts.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}

func TestGetPostsIntegration(t *testing.T) {
	// Initialize MongoDBPostRepository
	repo := repository.NewMongoDBPostRepository(testDBPosts)
	handlers.InitPostRepository(repo)

	// Create a new router and register the handler
	router := mux.NewRouter()
	router.HandleFunc("/posts", handlers.GetPosts).Methods("GET")

	// Seed some data into the database for testing
	seedPosts := []entity.Post{
		{Title: "Post 1", Content: "Content 1", CreatedAt: time.Now()},
		{Title: "Post 2", Content: "Content 2", CreatedAt: time.Now()},
	}
	for _, post := range seedPosts {
		_, err := repo.Save(&post)
		assert.Nil(t, err, "error seeding posts into database")
	}

	// Create a GET request to retrieve posts
	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code and response body
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	var posts []entity.Post
	err = json.Unmarshal(rr.Body.Bytes(), &posts)
	assert.Nil(t, err, "error unmarshalling response body")
	assert.NotNil(t, posts, "expected posts to be returned")
	assert.Equal(t, len(seedPosts), len(posts), "expected number of posts to match seeded data")
}

func TestAddPostIntegration(t *testing.T) {
	// Initialize MongoDBPostRepository
	repo := repository.NewMongoDBPostRepository(testDBPosts)
	handlers.InitPostRepository(repo)

	// Create a new router and register the handler
	router := mux.NewRouter()
	router.HandleFunc("/posts", handlers.AddPost).Methods("POST")

	// Create a sample post to add
	newPost := entity.Post{
		Title:   "Test Post",
		Content: "This is a test post.",
	}

	// Encode the post as JSON
	jsonPost, err := json.Marshal(newPost)
	if err != nil {
		t.Fatal(err)
	}

	// Create a POST request to add the post
	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonPost))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code and response body
	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
	var createdPost entity.Post
	err = json.Unmarshal(rr.Body.Bytes(), &createdPost)
	assert.Nil(t, err, "error unmarshalling response body")
	assert.Equal(t, newPost.Title, createdPost.Title, "expected title of created post to match")
	assert.Equal(t, newPost.Content, createdPost.Content, "expected content of created post to match")
}
