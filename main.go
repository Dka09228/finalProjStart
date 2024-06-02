package main

import (
	"finalProjStart/handlers"
	"finalProjStart/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	repo repository.PostRepository
)

func init() {
	var err error
	repo, err = repository.NewPostRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/api/posts", handlers.GetPosts(repo)).Methods("GET")
	router.HandleFunc("/api/posts/{id}", handlers.GetPostByID(repo)).Methods("GET")
	router.HandleFunc("/api/posts", handlers.AddPost(repo)).Methods("POST")
	router.HandleFunc("/api/posts/{id}", handlers.UpdatePost(repo)).Methods("PUT")
	router.HandleFunc("/api/posts/{id}", handlers.DeletePost(repo)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("Starting server...")
	handleRequests()
}
