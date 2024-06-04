package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"strings"
)

type TeamStats struct {
	Rank           int
	Team           string
	Games          int
	Wins           int
	Ties           int
	Losses         int
	GoalsFor       int
	GoalsAgainst   int
	GoalDiff       string
	Points         int
	PointsAvg      float64
	XGFor          float64
	XGAgainst      float64
	XGDiff         float64
	XGDiffPer90    float64
	AttendancePerG int
	TopTeamScorers string
	TopKeeper      string
	Notes          string
}

func main() {
	// Connect to MongoDB
	client, err := ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// URL of the page to scrape
	url := "https://fbref.com/en/comps/9/Premier-League-Stats"

	// Fetch the page
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	// Find the table
	doc.Find("table#results2023-202491_overall tbody tr").Each(func(i int, s *goquery.Selection) {
		teamStats := TeamStats{}

		// Scraping the data
		teamStats.Rank = i + 1
		teamStats.Team = strings.TrimSpace(s.Find("td[data-stat='team']").Text())
		teamStats.Games, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='games']").Text()))
		teamStats.Wins, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='wins']").Text()))
		teamStats.Ties, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='ties']").Text()))
		teamStats.Losses, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='losses']").Text()))
		teamStats.GoalsFor, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='goals_for']").Text()))
		teamStats.GoalsAgainst, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='goals_against']").Text()))
		teamStats.GoalDiff = strings.TrimSpace(s.Find("td[data-stat='goal_diff']").Text())
		teamStats.Points, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='points']").Text()))
		teamStats.PointsAvg, _ = strconv.ParseFloat(strings.TrimSpace(s.Find("td[data-stat='points_avg']").Text()), 64)
		// Additional stats (xG, attendance, top scorers, top keeper, notes)
		teamStats.XGFor, _ = strconv.ParseFloat(strings.TrimSpace(s.Find("td[data-stat='xg_for']").Text()), 64)
		teamStats.XGAgainst, _ = strconv.ParseFloat(strings.TrimSpace(s.Find("td[data-stat='xg_against']").Text()), 64)
		teamStats.XGDiff, _ = strconv.ParseFloat(strings.TrimSpace(s.Find("td[data-stat='xg_diff']").Text()), 64)
		teamStats.XGDiffPer90, _ = strconv.ParseFloat(strings.TrimSpace(s.Find("td[data-stat='xg_diff_per90']").Text()), 64)
		teamStats.AttendancePerG, _ = strconv.Atoi(strings.TrimSpace(s.Find("td[data-stat='attendance_per_g']").Text()))
		teamStats.TopTeamScorers = strings.TrimSpace(s.Find("td[data-stat='top_team_scorers']").Text())
		teamStats.TopKeeper = strings.TrimSpace(s.Find("td[data-stat='top_keeper']").Text())
		teamStats.Notes = strings.TrimSpace(s.Find("td[data-stat='notes']").Text())

		// Insert into MongoDB
		err := InsertTeamStats(client, teamStats)
		if err != nil {
			log.Printf("Failed to insert team stats for %s: %v", teamStats.Team, err)
		}
	})
}

// ConnectMongoDB creates and returns a connection to the MongoDB
func ConnectMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://alisheribraev03:m3Zi0rAnHCMJpm0f@alish.ahqrfiy.mongodb.net/?retryWrites=true&w=majority&appName=Alish")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

// InsertTeamStats inserts the given team stats into the MongoDB collection
func InsertTeamStats(client *mongo.Client, teamStats TeamStats) error {
	collection := client.Database("golang-test").Collection("premleague-test")
	_, err := collection.InsertOne(context.Background(), teamStats)
	if err != nil {
		return err
	}
	return nil
}
