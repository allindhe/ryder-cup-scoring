package main

import (
	"log"
	"net/http"
)

func main() {
	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/match/", getMatchHandler)
	http.HandleFunc("/totalScore/", getTotalScore)
	http.HandleFunc("/clear/", getClearMatches)

	log.Print("Listening on localhost:3000...")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
