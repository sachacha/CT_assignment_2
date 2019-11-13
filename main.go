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
	http.HandleFunc("/repocheck/v1/status", APIs.HandlerStatus)
	fmt.Println("Listening on port " + port)
	
	var START_TIME = time.Now()

	for {
		if (START_TIME.Second() > 10) {
			fmt.Println("et de 10 !")
		}
	}
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
