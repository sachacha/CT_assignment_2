package main

import (
	"CT_assignment_2/APIs"
	"fmt"
	"log"
	"net/http"
	"os"
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
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
