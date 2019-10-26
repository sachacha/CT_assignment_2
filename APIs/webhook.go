package APIs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	//"firebase.google.com/go/auth"
	"golang.org/x/net/context"

	"google.golang.org/api/iterator"
	//"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Registration struct {
	Event string `json:"event"`
	Url   string `json:"url"`
}

type RegistrationAnswer struct {
	ID string `json:"id"`
}

type RegisteredWebhook struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Time  string `json:"time"`
}

func HandlerWebhook(w http.ResponseWriter, r *http.Request) {

	// connection to the DB
	ctx := context.Background()

	sa := option.WithCredentialsFile("/home/sacha/Downloads/ctassignment2-firebase.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	switch r.Method {
	case http.MethodPost:
		var registration Registration
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("Decoding: " + err.Error())
			return
		}

		fmt.Println("Adding to db ...")

		ID, _, err := client.Collection("webhooks").Add(ctx, map[string]interface{}{
			"event": registration.Event,
			"url":   registration.Url,
		})
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}

		http.Header.Add(w.Header(), "content-type", "application/json")

		var registrationAnswer = RegistrationAnswer{}
		registrationAnswer.ID = ID.ID

		json.NewEncoder(w).Encode(registrationAnswer)

		return

	case http.MethodGet:
		http.Header.Add(w.Header(), "content-type", "application/json")

		var registeredWebhooks = []RegisteredWebhook{}

		iter := client.Collection("webhooks").Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			var registeredWebhook = RegisteredWebhook{}
			registeredWebhook.ID = doc.Ref.ID
			registeredWebhook.Event = doc.Data()["event"].(string)
			registeredWebhook.Time = doc.CreateTime.String()

			registeredWebhooks = append(registeredWebhooks, registeredWebhook)
		}

		json.NewEncoder(w).Encode(registeredWebhooks)

		return

	default:
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
		return
	}
}
