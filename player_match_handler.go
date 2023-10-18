package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Score struct {
	Team1 float32
	Team2 float32
}

type Match struct {
	// score: 0 = not played, 1 = team 1, 2 = draw, 3 = team 2
	score  [9]int
	result float32
}

type updateMatchStruct struct {
	Player1 string
	Player2 string
	Score   [9]int
}

var totalScore = Score{}

var matches map[string]*Match
var team1 = [3]string{"tom", "alex", "rasmus"}
var team2 = [3]string{"sebbe", "emil", "jean"}

func init() {
	clearMatches()
}

func clearMatches() {
	totalScore.Team1 = 0
	totalScore.Team2 = 0

	matches = make(map[string]*Match)

	for _, id1 := range team1 {
		for _, id2 := range team2 {
			newMatch := Match{}
			matches[id1+id2] = &newMatch
		}
	}

	// Create a special one for best ball
	newMatch := Match{}
	matches["bestball"] = &newMatch
}

func getClearMatches(w http.ResponseWriter, r *http.Request) {
	clearMatches()
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Data cleared")
}

func getTotalScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(totalScore)
}

func getMatchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query()
		player1 := q.Get("player1")
		player2 := q.Get("player2")

		if player1 == "" || player2 == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pmatch, ok := matches[player1+player2]
		if ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(pmatch.score)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		// Create a new record.
	case http.MethodPut:
		// Update an existing record.
		decoder := json.NewDecoder(r.Body)
		var t updateMatchStruct
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println("Invalid put format")
			w.WriteHeader(http.StatusInternalServerError)
		}

		pMatch, ok := matches[t.Player1+t.Player2]
		if ok {
			if len(t.Score) == 9 {
				for i, val := range t.Score {
					pMatch.score[i] = val
				}
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		updateScore(pMatch)
		updateTotalScore()

	case http.MethodDelete:
		// Remove the record.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func updateScore(pMatch *Match) {
	localScore := 0
	for _, val := range pMatch.score {
		if val == 0 {
			// Only count finished matches
			return
		}
		localScore += val - 2
	}

	if localScore == 0 {
		pMatch.result = 2
	} else if localScore < 0 {
		pMatch.result = 1
	} else {
		pMatch.result = 3
	}
}

func updateTotalScore() {
	localScore1 := float32(0)
	localScore2 := float32(0)

	for _, match := range matches {
		switch match.result {
		case 1:
			localScore1++
		case 2:
			localScore1 += 0.5
			localScore2 += 0.5
		case 3:
			localScore2++
		}
	}

	totalScore.Team1 = localScore1
	totalScore.Team2 = localScore2
}
