package main

import (
	"finalProjStart/db"
	"finalProjStart/handlers"
	"finalProjStart/jsonlog"
	"finalProjStart/middleware"
	"finalProjStart/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	logger   *jsonlog.Logger
)

func init() {
	logger = jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	client, err := db.ConnectMongoDB()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	database := client.Database(db.DatabaseName)

	postRepo = repository.NewMongoDBPostRepository(database)
	userRepo = repository.NewMongoDBUserRepository(database)

	handlers.InitUserRepository(userRepo)
	handlers.InitPostRepository(postRepo)
}

func handleRequests() {
	router := mux.NewRouter()
	router.Use(middleware.JSONMiddleware)
	router.Use(middleware.RecoverPanic)

	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutUser).Methods("POST") // Endpoint for logout

	router.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.GetPosts, "user", "admin")).Methods("GET")

	router.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handlers.GetPostByID, "user", "admin")).Methods("GET")
	router.HandleFunc("/api/posts", middleware.AuthMiddleware(handlers.AddPost, "admin")).Methods("POST")
	router.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handlers.UpdatePost, "admin")).Methods("PUT")
	router.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handlers.DeletePost, "admin")).Methods("DELETE")

	logger.PrintInfo("Starting server...", nil)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	handleRequests()
}
