package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"strings"
)

type TeamStats struct {
	Rank           int     `bson:"rank"`
	Team           string  `bson:"team"`
	Games          int     `bson:"games"`
	Wins           int     `bson:"wins"`
	Ties           int     `bson:"ties"`
	Losses         int     `bson:"losses"`
	GoalsFor       int     `bson:"goals_for"`
	GoalsAgainst   int     `bson:"goals_against"`
	GoalDiff       string  `bson:"goal_diff"`
	Points         int     `bson:"points"`
	PointsAvg      float64 `bson:"points_avg"`
	XGFor          float64 `bson:"xg_for"`
	XGAgainst      float64 `bson:"xg_against"`
	XGDiff         string  `bson:"xg_diff"`
	XGDiffPer90    string  `bson:"xg_diff_per90"`
	AttendancePerG int     `bson:"attendance_per_g"`
	TopTeamScorers string  `bson:"top_team_scorers"`
	TopKeeper      string  `bson:"top_keeper"`
	Notes          string  `bson:"notes"`
}

func main() {
	// Connect to MongoDB
	client, err := ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// URL of the page to scrape
	url := "https://fbref.com/en/comps/11/Serie-A-Stats"

	// Fetch the page
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	// Find the table
	doc.Find("table#results2023-2024111_overall tbody tr").Each(func(i int, s *goquery.Selection) {
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
		teamStats.XGDiff = strings.TrimSpace(s.Find("td[data-stat='xg_diff']").Text())
		teamStats.XGDiffPer90 = strings.TrimSpace(s.Find("td[data-stat='xg_diff_per90']").Text())
		attendanceStr := strings.ReplaceAll(strings.TrimSpace(s.Find("td[data-stat='attendance_per_g']").Text()), ",", "")
		teamStats.AttendancePerG, _ = strconv.Atoi(attendanceStr)
		teamStats.TopTeamScorers = strings.TrimSpace(s.Find("td[data-stat='top_team_scorers']").Text())
		teamStats.TopKeeper = strings.TrimSpace(s.Find("td[data-stat='top_keeper']").Text())
		teamStats.Notes = strings.TrimSpace(s.Find("td[data-stat='notes']").Text())

		// Log scraped data for debugging
		log.Printf("Scraped data for team %s: %+v", teamStats.Team, teamStats)

		// Insert into MongoDB
		err := InsertTeamStats(client, teamStats)
		if err != nil {
			log.Printf("Failed to insert team stats for %s: %v", teamStats.Team, err)
		} else {
			log.Printf("Successfully inserted/updated team stats for %s", teamStats.Team)
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
	collection := client.Database("golang-test").Collection("seriea-test")
	// Define the filter to find the document to update
	filter := bson.M{"team": teamStats.Team}

	// Define the update operation
	update := bson.D{
		{"$set", bson.D{
			{"rank", teamStats.Rank},
			{"team", teamStats.Team},
			{"games", teamStats.Games},
			{"wins", teamStats.Wins},
			{"ties", teamStats.Ties},
			{"losses", teamStats.Losses},
			{"goals_for", teamStats.GoalsFor},
			{"goals_against", teamStats.GoalsAgainst},
			{"goal_diff", teamStats.GoalDiff},
			{"points", teamStats.Points},
			{"points_avg", teamStats.PointsAvg},
			{"xg_for", teamStats.XGFor},
			{"xg_against", teamStats.XGAgainst},
			{"xg_diff", teamStats.XGDiff},
			{"xg_diff_per90", teamStats.XGDiffPer90},
			{"attendance_per_g", teamStats.AttendancePerG},
			{"top_team_scorers", teamStats.TopTeamScorers},
			{"top_keeper", teamStats.TopKeeper},
			{"notes", teamStats.Notes},
		}},
	}

	// Perform the upsert operation
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}
