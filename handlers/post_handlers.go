package handlers

import (
	"encoding/json"
	"finalProjStart/entity"
	"finalProjStart/error"
	"finalProjStart/jsonlog"
	"finalProjStart/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var (
	postRepo   repository.PostRepository
	postLogger *jsonlog.Logger // Use a unique name to avoid redeclaration
)

func InitPostRepository(repo repository.PostRepository) {
	postRepo = repo
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postRepo.FindAll()
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "retrieving posts"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to retrieve posts"))
		return
	}
	if posts == nil {
		postLogger.PrintInfo("No posts found", nil)
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusNotFound, "No posts found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "parsing ID from request"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid ID"))
		return
	}

	post, err := postRepo.FindByID(id)
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "retrieving post by ID"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to retrieve post"))
		return
	}
	if post == nil {
		postLogger.PrintInfo("Post not found", map[string]string{"postID": id.Hex()})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusNotFound, "Post not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "decoding request payload"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	_, err = postRepo.Save(&post)
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "saving post"})
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
		postLogger.PrintError(err, map[string]string{"context": "parsing ID from request"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid ID"))
		return
	}

	var updatedPost entity.Post
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "decoding request payload"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}
	updatedPost.ID = id

	err = postRepo.Update(&updatedPost)
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "updating post"})
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
		postLogger.PrintError(err, map[string]string{"context": "parsing ID from request"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusBadRequest, "Invalid ID"))
		return
	}

	err = postRepo.Delete(id)
	if err != nil {
		postLogger.PrintError(err, map[string]string{"context": "deleting post"})
		errors.SendErrorResponse(w, errors.NewAPIError(http.StatusInternalServerError, "Failed to delete post"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
