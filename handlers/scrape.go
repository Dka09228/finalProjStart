package handlers

import (
	"encoding/json"
	"finalProjStart/jsonlog"
	"finalProjStart/repository"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type ScrapeRequest struct {
	URL        string `json:"url"`
	TableID    string `json:"table_id"`
	Collection string `json:"collection"`
}

var (
	scrapeLogger *jsonlog.Logger
	scrapeRepo   repository.ScrapeRepository
)

func InitScrapeRepository(repo repository.ScrapeRepository) {
	scrapeRepo = repo
}

func ScrapeData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := scrapeRepo.ScrapeAndStoreData(req.URL, req.TableID, req.Collection); err != nil {
		http.Error(w, "Failed to scrape and store data: "+err.Error(), http.StatusInternalServerError)
		scrapeLogger.PrintError(err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data scraped and stored successfully")
	scrapeLogger.PrintInfo("Data scraped and stored successfully", nil)
}

func GetTeamStatsByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	collectionName := r.URL.Query().Get("collection")

	teamStats, err := scrapeRepo.GetTeamStatsByID(id, collectionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teamStats)
}

func GetTeamStatsByTeamName(w http.ResponseWriter, r *http.Request) {
	teamName := mux.Vars(r)["team"]
	collectionName := r.URL.Query().Get("collection")

	teamStats, err := scrapeRepo.GetTeamStatsByTeamName(teamName, collectionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teamStats)
}

func GetLeagueByCollectionName(w http.ResponseWriter, r *http.Request) {
	collectionName := r.URL.Query().Get("collection")

	league, err := scrapeRepo.GetLeagueByCollectionName(collectionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(league)
}

func DeleteLeague(w http.ResponseWriter, r *http.Request) {
	collectionName := r.URL.Query().Get("collection")

	err := scrapeRepo.DeleteLeague(collectionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
