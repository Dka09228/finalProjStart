package main

import (
	"log"
	"net/http"

	"finalProjStart/repository"
	"github.com/gorilla/mux"
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
	router.HandleFunc("/api/posts", getPosts).Methods("GET")
	router.HandleFunc("/api/posts", addPost).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("Starting server...")
	handleRequests()
}
