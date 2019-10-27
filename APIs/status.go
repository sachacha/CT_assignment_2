package APIs

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var START_TIME = time.Now()

type StatusAnswer struct {
	Gitlab   int     `json:"gitlab"`
	Database int     `json:"database"`
	Uptime   float64 `json:"uptime"`
	Version  string  `json:"version"`
}

func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	var parameters []string
	WebhookChecking(w, "languages", parameters)

	// check gitlab availability
	respGit, errGit := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")

	if errGit != nil {
		http.Error(w, errGit.Error(), http.StatusBadRequest)
		return
	}

	// check the database availability
	dbAvailability := 200

	ctx := context.Background()

	sa := option.WithCredentialsFile("/home/sacha/Downloads/ctassignment2-firebase.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
		dbAvailability = 404
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		dbAvailability = 404
	}

	defer client.Close()

	uptime := time.Since(START_TIME).Seconds()

	statusAnswer := StatusAnswer{respGit.StatusCode, dbAvailability, uptime, "v1"}

	json.NewEncoder(w).Encode(statusAnswer)
}
