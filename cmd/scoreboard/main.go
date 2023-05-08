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
	"time"

	pkg "github.com/dapr-volleyball-demo/pkg"

	client "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	dapr "github.com/dapr/go-sdk/service/http"
)

const stateStoreComponentName = "statestore"

var sub = &common.Subscription{
	PubsubName: "gamepubsub",
	Topic:      "game",
	Route:      "/save-game",
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3002"
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
	log.Printf("starting scoreboard service")
	err = s.Start()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}

// curl -X POST http://localhost:3002/echo -H "Content-Type: application/json" -d '{"message": "hello world"}'
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

// eventHandler receives data on the game topic and saves state on game point of 25 or higher for either team.
func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] Subscriber received data %v and rawdata %v\n", currentTime, e.Data, e.RawData)
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the incoming score message
	var game pkg.Game
	err = json.Unmarshal(e.RawData, &game)
	if err != nil {
		log.Fatalf("error unmarshalling into game %v", err)
	}

	// Save state into the state store if game point or higher (ie point 25 or higher)
	if game.Team1Score > 25 || game.Team2Score > 25 {
		key := "game_" + strconv.Itoa(game.ID)
		err = client.SaveState(context.Background(), stateStoreComponentName, key, e.RawData, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[%s] Saved game score: %s\n", currentTime, string(e.RawData))
	}

	return false, nil
}

// curl -X POST http://localhost:3002/scoreboard -H "Content-Type: application/json" -d '{"id": 0}'
func getGameScoreboardHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	if in == nil {
		err = errors.New("invocation parameter required")
		return
	}
	log.Printf(
		"[%s] echo - ContentType:%s, Verb:%s, QueryString:%s, %s",
		currentTime, in.ContentType, in.Verb, in.QueryString, in.Data,
	)

	var gameReq pkg.GameRequest
	err = json.Unmarshal(in.Data, &gameReq)
	if err != nil {
		log.Printf("error unmarshalling into gameReq")
		return nil, err
	}

	// Get the state from the state store using the game ID
	client, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	key := "game_" + strconv.Itoa(gameReq.ID)
	item, err := client.GetState(context.Background(), stateStoreComponentName, key, nil)
	if err != nil {
		log.Printf("error getting state for id %d", &gameReq.ID)
		return
	}
	log.Printf("[%s] retrieved state for game: %s", currentTime, string(item.Value))

	out = &common.Content{
		Data:        item.Value,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}
