// handlers/post_handlers.go
package handlers

import (
	"encoding/json"
	"finalProjStart/entity"
	"finalProjStart/error"
	"finalProjStart/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var postRepo repository.PostRepository

func InitPostRepository(repo repository.PostRepository) {
	postRepo = repo
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postRepo.FindAll()
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to retrieve posts"))
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid ID"))
		return
	}

	post, err := postRepo.FindByID(id)
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to retrieve post"))
		return
	}
	if post == nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusNotFound, "Post not found"))
		return
	}
	json.NewEncoder(w).Encode(post)
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	_, err = postRepo.Save(&post)
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to save post"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid ID"))
		return
	}

	var updatedPost entity.Post
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}
	updatedPost.ID = id

	err = postRepo.Update(&updatedPost)
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to update post"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid ID"))
		return
	}

	err = postRepo.Delete(id)
	if err != nil {
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to delete post"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
