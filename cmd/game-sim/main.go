package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	pkg "github.com/dapr-volleyball-demo/pkg"
	dapr "github.com/dapr/go-sdk/client"
)

const (
	maxPoints     = 25
	minPointsDiff = 2

	pubsubComponentName = "gamepubsub"
	pubsubTopic         = "game"
)

func main() {
	// Create a new client for Dapr using the SDK
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Publish events using Dapr pubsub
	// simulate 5 games to play
	for i := 0; i < 5; i++ {
		var game pkg.Game
		game.ID = i
		game.Team1Name = "team" + strconv.Itoa(i)
		game.Team2Name = "team" + strconv.Itoa(i+1)
		for {
			if game.Team1Score >= maxPoints && game.Team1Score-game.Team2Score >= minPointsDiff {
				log.Printf("team 1 wins: %+v", game)
				break
			}

			if game.Team2Score >= maxPoints && game.Team2Score-game.Team1Score >= minPointsDiff {
				log.Printf("team 2 wins: %+v", game)
				break
			}

			// Simulate the game by randomly incrementing one team's score.
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(2) == 0 {
				game.Team1Score++
			} else {
				game.Team2Score++
			}
			game.Round++

			err = client.PublishEvent(context.Background(), pubsubComponentName, pubsubTopic, game)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Published data: %#v\n", game)

			time.Sleep(1000)
		}
	}
}
