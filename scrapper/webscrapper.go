package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func main() {
	// URL of the page to scrape
	url := "https://fbref.com/en/comps/9/Premier-League-Stats"

	// Fetch the page
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	// Find the table
	doc.Find("table#results2023-202491_overall tbody tr").Each(func(i int, s *goquery.Selection) {
		rank := s.Find("th[data-stat='rank']").Text()
		team := s.Find("td[data-stat='team']").Text()
		games := s.Find("td[data-stat='games']").Text()
		wins := s.Find("td[data-stat='wins']").Text()
		ties := s.Find("td[data-stat='ties']").Text()
		losses := s.Find("td[data-stat='losses']").Text()
		goalsFor := s.Find("td[data-stat='goals_for']").Text()
		goalsAgainst := s.Find("td[data-stat='goals_against']").Text()
		goalDiff := s.Find("td[data-stat='goal_diff']").Text()
		points := s.Find("td[data-stat='points']").Text()
		pointsAvg := s.Find("td[data-stat='points_avg']").Text()

		// Additional stats (xG, attendance, top scorers, top keeper, notes)
		xgFor := s.Find("td[data-stat='xg_for']").Text()
		xgAgainst := s.Find("td[data-stat='xg_against']").Text()
		xgDiff := s.Find("td[data-stat='xg_diff']").Text()
		xgDiffPer90 := s.Find("td[data-stat='xg_diff_per90']").Text()
		attendancePerG := s.Find("td[data-stat='attendance_per_g']").Text()
		topTeamScorers := s.Find("td[data-stat='top_team_scorers']").Text()
		topKeeper := s.Find("td[data-stat='top_keeper']").Text()
		notes := s.Find("td[data-stat='notes']").Text()

		// Clean up the text
		rank = "Rank: " + strings.TrimSpace(rank)
		team = "Team: " + strings.TrimSpace(team)
		games = "G: " + strings.TrimSpace(games)
		wins = "W: " + strings.TrimSpace(wins)
		ties = "D: " + strings.TrimSpace(ties)
		losses = "L: " + strings.TrimSpace(losses)
		goalsFor = "GF: " + strings.TrimSpace(goalsFor)
		goalsAgainst = "GA: " + strings.TrimSpace(goalsAgainst)
		goalDiff = "GD: " + strings.TrimSpace(goalDiff)
		points = "Pts: " + strings.TrimSpace(points)
		pointsAvg = "Pts/MP: " + strings.TrimSpace(pointsAvg)
		xgFor = "xG: " + strings.TrimSpace(xgFor)
		xgAgainst = "xGA: " + strings.TrimSpace(xgAgainst)
		xgDiff = "xGD: " + strings.TrimSpace(xgDiff)
		xgDiffPer90 = "xGD/90: " + strings.TrimSpace(xgDiffPer90)
		attendancePerG = "Attendance: " + strings.TrimSpace(attendancePerG)
		topTeamScorers = "Top Scorer: " + strings.TrimSpace(topTeamScorers)
		topKeeper = "Top Goalkeeper: " + strings.TrimSpace(topKeeper)
		notes = strings.TrimSpace(notes)

		// Print the result
		fmt.Printf("%s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n| %s | %s | %s\n",
			rank, team, games, wins, ties, losses, goalsFor, goalsAgainst, goalDiff, points, pointsAvg, xgFor, xgAgainst, xgDiff, xgDiffPer90, attendancePerG, topTeamScorers, topKeeper, notes)
	})
}
