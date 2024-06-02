package handlers

import (
	"encoding/json"
	"finalProjStart/entity"
	"finalProjStart/repository"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func GetPosts(repo repository.PostRepository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := repo.FindAll()
		if err != nil {
			http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func GetPostByID(repo repository.PostRepository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		post, err := repo.FindByID(id)
		if err != nil {
			http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
			return
		}
		if post == nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func AddPost(repo repository.PostRepository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post entity.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		newPost, err := repo.Save(&post)
		if err != nil {
			http.Error(w, "Failed to add post", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newPost)
	}
}

func UpdatePost(repo repository.PostRepository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var post entity.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		updatedPost, err := repo.Update(id, &post)
		if err != nil {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}
		if updatedPost == nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedPost)
	}
}

func DeletePost(repo repository.PostRepository) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = repo.Delete(id)
		if err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
