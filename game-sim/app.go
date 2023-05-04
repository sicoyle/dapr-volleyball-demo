package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	maxPoints     = 25
	minPointsDiff = 2
)

const (
	pubsubComponentName = "gamepubsub"
	pubsubTopic         = "game"
)

type Scoreboard struct {
	Round      int `json:"round"` // change to be sets
	Team1Score int `json:"team1Score"`
	Team2Score int `json:"team2Score"`
}

func main() {
	// Create a new client for Dapr using the SDK
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Publish events using Dapr pubsub
	var score Scoreboard
	for {
		if score.Team1Score >= maxPoints && score.Team1Score-score.Team2Score >= minPointsDiff {
			log.Printf("team 1 wins: %+v", score)
			return
		}

		if score.Team2Score >= maxPoints && score.Team2Score-score.Team1Score >= minPointsDiff {
			log.Printf("team 2 wins: %+v", score)
			return
		}

		// Simulate the game by randomly incrementing one team's score.
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(2) == 0 {
			score.Team1Score++
		} else {
			score.Team2Score++
		}
		score.Round++

		// score := `{"round":` + strconv.Itoa(i) + `,"team1Score":` + strconv.Itoa(i) + `,"team2Score":` + strconv.Itoa(i) + `}`

		// Publish the score update to the "game" topic
		data, err := json.Marshal(score)
		if err != nil {
			log.Printf("failed to marshal score message: %v", err)
			return
		}

		err = client.PublishEvent(context.Background(), pubsubComponentName, pubsubTopic, data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Published data: %#v\n", score)

		time.Sleep(1000)
	}
}
