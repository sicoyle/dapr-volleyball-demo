package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	client "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	dapr "github.com/dapr/go-sdk/service/http"
)

const stateStoreComponentName = "statestore"

// Score holds the score of the team who just made a point
// Scoreboard holds the current score of the volleyball game.
type Scoreboard struct {
	Round      int `json:"round"`
	Team1Score int `json:"team1Score"`
	Team2Score int `json:"team2Score"`
}

var sub = &common.Subscription{
	PubsubName: "gamepubsub",
	Topic:      "game",
	Route:      "/save-game",
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "6005"
	}

	// Create the new server on appPort and add a topic listener
	s := dapr.NewService(":" + appPort)

	err := s.AddTopicEventHandler(sub, eventHandler)
	if err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	// add a service to service invocation handler
	if err := s.AddServiceInvocationHandler("/echo", echoHandler); err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}
	if err := s.AddServiceInvocationHandler("/scoreboard", getGameScoreboardHandler); err != nil {
		log.Fatalf("error adding invocation handler for scoreboard: %v", err)
	}

	// Start the server
	err = s.Start()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}

// curl -X POST http://localhost:6006/echo -H "Content-Type: application/json" -d '{"message": "hello world"}'
func echoHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	if in == nil {
		err = errors.New("invocation parameter required")
		return
	}
	log.Printf(
		"echo - ContentType:%s, Verb:%s, QueryString:%s, %s",
		in.ContentType, in.Verb, in.QueryString, in.Data,
	)
	out = &common.Content{
		Data:        in.Data,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	fmt.Printf("Subscriber received data %v and rawdata %v\n", e.Data, e.RawData)
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the incoming score message
	var score Scoreboard
	// Save state into the state store
	err = client.SaveState(context.Background(), stateStoreComponentName, strconv.Itoa(score.Round), e.RawData, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved Score:", string(e.RawData))

	return false, nil
}

// curl -X POST http://localhost:6006/scoreboard -H "Content-Type: application/json" -d '{"round": 45}'
func getGameScoreboardHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	if in == nil {
		err = errors.New("invocation parameter required")
		return
	}
	log.Printf(
		"echo - ContentType:%s, Verb:%s, QueryString:%s, %s",
		in.ContentType, in.Verb, in.QueryString, in.Data,
	)

	var scoreboardReq Scoreboard
	err = json.Unmarshal(in.Data, &scoreboardReq)
	if err != nil {
		log.Printf("error unmarshalling into scoreboardReq")
		return nil, err
	}

	// Get the state from the state store using the game UID
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	item, err := client.GetState(context.Background(), stateStoreComponentName, "0", nil)
	if err != nil {
		log.Printf("error getting state for id %v", &scoreboardReq.Round)
		return
	}
	log.Printf("string value %s", string(item.Value))

	out = &common.Content{
		Data:        item.Value,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}

func getGameScoreboard(w http.ResponseWriter, r *http.Request) {
	// Get the game UID from the URL path
	gameUID := r.URL.Query().Get("gameUID")

	// Get the state from the state store using the game UID
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	item, err := client.GetState(context.Background(), stateStoreComponentName, gameUID, nil)
	if err != nil {
		http.Error(w, "Error getting state from store", http.StatusInternalServerError)
		return
	}

	// Unmarshal the state and return it as the response
	var scoreboard Scoreboard
	err = json.Unmarshal(item.Value, &scoreboard)
	if err != nil {
		http.Error(w, "Error unmarshaling state", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(scoreboard)
}
