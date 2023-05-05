package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/gorilla/mux"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	// TODO cleanup, but workaround bc kept getting err without this setup:
	// "error invoking rpc error: code = Canceled desc = grpc: the client connection is closing"
	client, cancel = newDaprClient()
)

func newDaprClient() (dapr.Client, func()) {
	client, err := dapr.NewClient()
	if err != nil {
		// TODO handle error
	}
	return client, func() {
		defer client.Close()
	}
}

func main() {
	defer cancel()
	router := mux.NewRouter()
	router.HandleFunc("/scoreboard", scoreboardHandler)
	srv := daprd.NewServiceWithMux(":8080", router)

	// Start the Dapr service
	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		log.Printf("error: %v", err)
	}
}

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	content := &dapr.DataContent{
		Data:        []byte(`{"round":45,"team1Score":0,"team2Score":0}`),
		ContentType: "application/json",
	}

	// invoke the service
	resp, err := client.InvokeMethodWithContent(context.Background(), "scoreboard", "scoreboard", "POST", content)
	if err != nil {
		log.Printf("error invoking %v", err)
	}

	// process the response
	fmt.Println(string(resp))
	json.NewEncoder(w).Encode(string(resp))
}
