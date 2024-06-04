package handlers

import (
	"encoding/json"
	"finalProjStart/jsonlog"
	"finalProjStart/repository"
	"fmt"
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
