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
	log.Printf("starting scoreboard service")
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
	var game pkg.Game
	err = json.Unmarshal(e.RawData, &game)
	if err != nil {
		log.Fatalf("error unmarshalling into game %v", err.Error())
		// return nil, err
	}

	// Save state into the state store
	key := "game_" + strconv.Itoa(game.ID)
	err = client.SaveState(context.Background(), stateStoreComponentName, key, e.RawData, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved game score:", string(e.RawData))

	return false, nil
}

// curl -X POST http://localhost:6006/scoreboard -H "Content-Type: application/json" -d '{"id": 0}'
func getGameScoreboardHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	if in == nil {
		err = errors.New("invocation parameter required")
		return
	}
	log.Printf(
		"echo - ContentType:%s, Verb:%s, QueryString:%s, %s",
		in.ContentType, in.Verb, in.QueryString, in.Data,
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
	log.Printf("string value %s", string(item.Value))

	out = &common.Content{
		Data:        item.Value,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return
}
