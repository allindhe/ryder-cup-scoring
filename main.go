package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var client *firestore.Client
var ctx context.Context

func main() {
	// Setup database
	ctx = context.Background()
	sa := option.WithCredentialsFile("grindr-cup-07c591b59b67.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/match/", getMatchHandler)
	http.HandleFunc("/totalScore/", getTotalScore)
	http.HandleFunc("/clear/", getClearMatches)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
