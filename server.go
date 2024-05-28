package main

import (
	"encoding/json"
	"finalProjStart/entity"
	"finalProjStart/repository"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"log"
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

func getPosts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	posts, err := repo.FindAll()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"error":"Error Getting the Posts"}`))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(posts)
}

func addPost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var post entity.Post
	err := json.NewDecoder(request.Body).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"error":"Error Unmarshalling Data"}`))
		return
	}
	post.ID = rand.Int63()
	_, err = repo.Save(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"error":"Error Saving the Post"}`))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(post)
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
