package main

import (
	"finalProjStart/db"
	"finalProjStart/handlers"
	"finalProjStart/middleware"
	"finalProjStart/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var postRepo repository.PostRepository
var userRepo repository.UserRepository

func init() {
	client, err := db.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	database := client.Database(db.DatabaseName)

	// Initialize repositories
	postRepo = repository.NewMongoDBPostRepository(database)
	userRepo = repository.NewMongoDBUserRepository(database)

	// Initialize handlers with repositories
	handlers.InitUserRepository(userRepo)
	handlers.InitPostRepository(postRepo)
}

func handleRequests() {
	router := mux.NewRouter()

	// Routes for authentication and user management
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutUser).Methods("POST") // Endpoint for logout

	// Routes for posts management with authentication middleware
	router.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.GetPosts, "user")).Methods("GET")
	router.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handlers.GetPostByID, "user")).Methods("GET")
	router.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.AddPost, "admin")).Methods("POST")
	router.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handlers.UpdatePost, "admin")).Methods("PUT")
	router.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handlers.DeletePost, "admin")).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("Starting server...")
	handleRequests()
}
