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
	gameCount     = 100

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
	// simulate 100 games to play
	for i := 0; i < gameCount; i++ {
		var game pkg.Game
		game.ID = i
		game.Team1Name = "team" + strconv.Itoa(i)
		game.Team2Name = "team" + strconv.Itoa(i+1)
		for {
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			if game.Team1Score >= maxPoints && game.Team1Score-game.Team2Score >= minPointsDiff {
				log.Printf("[%s] team 1 wins: %+v", currentTime, game)
				break
			}

			if game.Team2Score >= maxPoints && game.Team2Score-game.Team1Score >= minPointsDiff {
				log.Printf("[%s] team 2 wins: %+v", currentTime, game)
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

			fmt.Printf("[%s] Published data: %#v\n", currentTime, game)

			time.Sleep(6 * time.Second)
		}
	}

	// Note: the following is added so the container keeps running for the demo.
	stop := make(chan struct{})
	<-stop // block the main goroutine from exiting
}
