package APIs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	firebase "firebase.google.com/go"
	"golang.org/x/net/context"

	"google.golang.org/api/iterator"
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

var keyAddress = "C:\\Users\\user\\Downloads\\ctassignment2-firebase.json"

func HandlerWebhook(w http.ResponseWriter, r *http.Request) {

	// connection to the DB
	ctx := context.Background()

	sa := option.WithCredentialsFile(keyAddress)
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

func HandlerWebhookWithId(w http.ResponseWriter, r *http.Request) {

	// connection to the DB
	ctx := context.Background()

	sa := option.WithCredentialsFile(keyAddress)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	http.Header.Add(w.Header(), "content-type", "application/json")

	parts := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case http.MethodGet:

		var registeredWebhooks = map[string]RegisteredWebhook{}

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

			registeredWebhooks[doc.Ref.ID] = registeredWebhook
		}

		json.NewEncoder(w).Encode(registeredWebhooks[parts[4]])

		return

	case http.MethodDelete:

		_, err := client.Collection("webhooks").Doc(parts[4]).Delete(ctx)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return
	}
}

type InfoSend struct {
	Event  string   `json:"event"`
	Params []string `json:"params"`
	Time   string   `json:"time"`
}

type RegisteredWebhookWithUrl struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Url   string `json:"url"`
}

func WebhookChecking(w http.ResponseWriter, eventType string, parameters []string) {
	// create the payload we will send
	infoSend := InfoSend{eventType, parameters, time.Now().String()}

	// connection to the DB
	ctx := context.Background()

	sa := option.WithCredentialsFile(keyAddress)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	// get the webhooks in the database
	var registeredWebhooks = map[string]RegisteredWebhookWithUrl{}

	iter := client.Collection("webhooks").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		var registeredWebhook = RegisteredWebhookWithUrl{}
		registeredWebhook.ID = doc.Ref.ID
		registeredWebhook.Event = doc.Data()["event"].(string)
		registeredWebhook.Url = doc.Data()["url"].(string)

		registeredWebhooks[doc.Ref.ID] = registeredWebhook
	}

	// get those which are related to our eventType
	var ids []string

	for webhook := range registeredWebhooks {
		if registeredWebhooks[webhook].Event == eventType {
			ids = append(ids, registeredWebhooks[webhook].ID)
		}
	}

	// send them our get request with the payload
	lenIds := len(ids)

	for i := 0; i < lenIds; i++ {
		url := registeredWebhooks[ids[i]].Url

		requestBody, err := json.Marshal(infoSend)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			continue
		}

		payload := bytes.NewBuffer(requestBody)

		req, err1 := http.NewRequest(http.MethodGet, url, payload)

		if err1 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			continue
		}
		
		_,_=http.DefaultClient.Do(req)
	}
}
