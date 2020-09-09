package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
    "google.golang.org/api/option"
)


var app *firebase.App

func sendToTopic(ctx context.Context, topic string, data map[string]interface{}, notification *messaging.Notification) error {
    // https://github.com/firebase/firebase-admin-go/blob/cef91acd46f2fc5d0b3408d8154a0005db5bdb0b/snippets/messaging.go#L59
	// Obtain a messaging.Client from the App.
    if app == nil {
        panic("what the fuck")
    }
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
    mapString := make(map[string]string)

    for key, value := range  data{
        strValue := fmt.Sprintf("%v", value)
        mapString[key] = strValue
    }

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: mapString,
        Notification: notification,
		Topic: topic,
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := client.Send(ctx, message)
	if err == nil {
		fmt.Println(response)
	}
	return err
}


func sendToToken(ctx context.Context, registrationToken string, data map[string]interface{}, notification *messaging.Notification) error {
    // https://github.com/firebase/firebase-admin-go/blob/cef91acd46f2fc5d0b3408d8154a0005db5bdb0b/snippets/messaging.go#L27
	// Obtain a messaging.Client from the App.
    if app == nil {
        panic("what the fuck")
    }
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
    mapString := make(map[string]string)

    for key, value := range  data{
        strValue := fmt.Sprintf("%v", value)
        mapString[key] = strValue
    }

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: mapString,
        Notification: notification,
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err == nil{
		fmt.Println("Successfully sent message:", response)
	}else {
		fmt.Println("Error sending message:", err)
	}
    return err
	// [END send_to_token_golang]
}


func sendMulticast(ctx context.Context, registrationTokens []string, data map[string]interface{}, notification *messaging.Notification) error {
    // https://github.com/firebase/firebase-admin-go/blob/969e50e3996254cdb245d057bb2618fbd64ff425/snippets/messaging.go#L143
    // [START send_multicast]
    // Create a list containing up to 100 registration tokens.
    // This registration tokens come from the client FCM SDKs.
    // Obtain a messaging.Client from the App.
    client, err := app.Messaging(ctx)
    if err != nil {
        log.Fatalf("error getting Messaging client: %v\n", err)
    }
    mapString := make(map[string]string)
    for key, value := range  data{
        strValue := fmt.Sprintf("%v", value)
        mapString[key] = strValue
    }
    message := &messaging.MulticastMessage{
        Data: mapString,
        Notification: notification,
        Tokens: registrationTokens,
    }
    br, err := client.SendMulticast(context.Background(), message)
    fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
    // [END send_multicast]
    return err
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Code based on https://firebase.google.com/docs/cloud-messaging/send-message#go
    w.WriteHeader(http.StatusOK)
}

func SendToTopicHandler(w http.ResponseWriter, r *http.Request) {
    // Code based on https://firebase.google.com/docs/cloud-messaging/send-message#go
    var m map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&m)
    if err != nil{
        panic(err)
    }
    var notification *messaging.Notification
    if n, ok := m["notification"]; ok {
        nn := n.(map[string]interface{})
        notification = &messaging.Notification{
            Title: nn["title"].(string),
            Body: nn["body"].(string),
        }
    }
    err = sendToTopic(r.Context(), m["topic"].(string), m["data"].(map[string]interface{}), notification)
    if err != nil{
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte(err.Error()))
    }else{
        w.WriteHeader(http.StatusOK)
	}
}

func SendToTokenHandler(w http.ResponseWriter, r *http.Request) {
    // Code based on https://firebase.google.com/docs/cloud-messaging/send-message#go
    var m map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&m)
    if err != nil{
        panic(err)
    }
    var notification *messaging.Notification
    if n, ok := m["notification"]; ok {
        nn := n.(map[string]interface{})
        notification = &messaging.Notification{
            Title: nn["title"].(string),
            Body: nn["body"].(string),
        }
    }
    err = sendToToken(r.Context(), m["token"].(string), m["data"].(map[string]interface{}), notification)
    if err != nil{
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte(err.Error()))
    }else{
        w.WriteHeader(http.StatusOK)
	}
}

func SendToTokensHandler(w http.ResponseWriter, r *http.Request) {
    // Code based on https://firebase.google.com/docs/cloud-messaging/send-message#go
    var m map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&m)
    if err != nil{
        panic(err)
    }
    var notification *messaging.Notification
    if n, ok := m["notification"]; ok {
        nn := n.(map[string]interface{})
        notification = &messaging.Notification{
            Title: nn["title"].(string),
            Body: nn["body"].(string),
        }
    }
	mm := m["tokens"].([]interface{})
	tokens := make([]string, len(mm))
	for ind, tok := range mm {
		tokens[ind] = tok.(string)
	}
    err = sendMulticast(r.Context(), tokens, m["data"].(map[string]interface{}), notification)
    if err != nil{
        w.WriteHeader(http.StatusNotAcceptable)
        w.Write([]byte(err.Error()))
    }else{
        w.WriteHeader(http.StatusOK)
	}
}
func main() {
    // INIT APP
    // https://github.com/firebase/firebase-admin-go/blob/cef91acd46f2fc5d0b3408d8154a0005db5bdb0b/snippets/init.go#L33
    APP, err := firebase.NewApp(context.Background(), nil)
    if err != nil {
        fmt.Println("Tyring to use CREDENTIALS_FILE.json.")
        opt := option.WithCredentialsFile("CREDENTIALS_FILE.json")
        APP, err = firebase.NewApp(context.Background(), nil, opt)
        if err != nil {
            log.Fatalf("error initializing app: %v\n", err)
        }
    }
    app = APP
    fmt.Println("USAGE:\n")
    fmt.Println("\t"+`curl localhost:8080/sendToToken -d '{
	    "data": {"k": "v"},
	    "notification": {"title": "t", "body": "b"},
	    "token":"<YOUR TOKEN>"
            }'`)
    fmt.Println()
    fmt.Println("\t"+`curl localhost:8080/sendToTokens -d '{
	    "data": {"k": "v"},
	    "notification": {"title": "t", "body": "b"},
	    "tokens": ["<TOKEN-0>", "<TOKEN-1>", ...]
            }'`)
    fmt.Println()
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/sendToToken", SendToTokenHandler)
    r.HandleFunc("/sendToTopic", SendToTopicHandler)
    r.HandleFunc("/sendToTokens", SendToTokensHandler)
    port := "8080"
    fmt.Println("Starting Firebase Cloud Messaging Notification Generator server on port:", port)
    loggedRouter := handlers.LoggingHandler(os.Stdout, r)
    srv := &http.Server{
        Handler:      loggedRouter,
        Addr:         "127.0.0.1:"+port, // :8000",
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    fmt.Printf("%+v\n",srv.Addr)
    fmt.Println()
    log.Fatal(srv.ListenAndServe())
}
