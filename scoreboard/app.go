package main

import (
	"context"
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
	Route:      "/game", // TODO rename this to save-game
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

	// Start the server
	err = s.Start()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
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
