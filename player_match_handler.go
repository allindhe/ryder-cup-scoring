package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/api/iterator"
)

type Score struct {
	Team1 float32
	Team2 float32
}

type Match struct {
	// score: 0 = not played, 1 = team 1, 2 = draw, 3 = team 2
	Score  [9]int
	Result float32
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

func clearMatches() {
	_, err := client.Collection("score").Doc("score").Set(ctx, Score{})
	if err != nil {
		fmt.Printf("Err clearing score: %s", err)
	}

	// matches = make(map[string]*Match)

	for _, id1 := range team1 {
		for _, id2 := range team2 {
			_, err = client.Collection("matches").Doc(id1+id2).Set(ctx, Match{})
			if err != nil {
				fmt.Printf("Err clearing match: %s", err)
			}
		}
	}

	// Create a special one for best ball
	_, err = client.Collection("matches").Doc("bestball").Set(ctx, Match{})
	if err != nil {
		fmt.Printf("Err clearing match: %s", err)
	}
}

func getClearMatches(w http.ResponseWriter, r *http.Request) {
	clearMatches()
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Data cleared")
}

func getTotalScore(w http.ResponseWriter, r *http.Request) {
	dsnap, err := client.Collection("score").Doc("score").Get(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var totalScore Score
	dsnap.DataTo(&totalScore)

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

		dsnap, err := client.Collection("matches").Doc(player1 + player2).Get(ctx)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var match Match
		dsnap.DataTo(&match)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(match.Score)

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
			return
		}

		tmpMatch := Match{}
		if len(t.Score) == 9 {
			for i, val := range t.Score {
				tmpMatch.Score[i] = val
			}
		}
		updateResult(&tmpMatch)

		_, err = client.Collection("matches").Doc(t.Player1+t.Player2).Set(ctx, tmpMatch)
		if err != nil {
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		updateTotalScore()

	case http.MethodDelete:
		// Remove the record.
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func updateResult(pMatch *Match) {
	localScore := 0
	for _, val := range pMatch.Score {
		if val == 0 {
			// Only count finished matches
			return
		}
		localScore += val - 2
	}

	if localScore == 0 {
		pMatch.Result = 2
	} else if localScore < 0 {
		pMatch.Result = 1
	} else {
		pMatch.Result = 3
	}
}

func updateTotalScore() {
	score := Score{}

	// Get all matches
	iter := client.Collection("matches").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		match := Match{}
		doc.DataTo(&match)

		switch match.Result {
		case 1:
			score.Team1++
		case 2:
			score.Team1 += 0.5
			score.Team2 += 0.5
		case 3:
			score.Team2++
		}
	}

	// Write total score to db
	_, err := client.Collection("score").Doc("score").Set(ctx, score)
	if err != nil {
		fmt.Printf("Err writing score: %s", err)
	}
}
