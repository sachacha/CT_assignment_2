package main

import (
	"CT_assignment_2/APIs"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/repocheck/v1/commits", APIs.HandlerCommits)
	http.HandleFunc("/repocheck/v1/languages", APIs.HandlerLanguages)
	http.HandleFunc("/repocheck/v1/webhooks/", APIs.HandlerWebhookWithId)
	http.HandleFunc("/repocheck/v1/webhooks", APIs.HandlerWebhook)
	http.HandleFunc("/repocheck/v1/status", APIs.HandlerStatus)
	fmt.Println("Listening on port " + port)

	var ticker = time.Now()

	for {
		if time.Since(ticker).Seconds() >= 10 {
			fmt.Println("Done !!")
			ticker = time.Now()
		}
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
