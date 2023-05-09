# Dapr Volleyball Demo

## Volleyball Game Simulator

Volleyball Game Simulator simulates a volleyball game scenario where there are 100 volleyball games.
It randomly adds a point to one of two teams during the game until one team wins by 2.
Game point in volleyball is set to 25, but there is no cap in our simulation.
As the game continues, it sends score updates onto the `game` topic of the `gamepubsub` pubsub.

```
cd cmd/game-sim

dapr run \
--app-id game-sim \
--app-protocol http \
--dapr-http-port 3500 \
--resources-path ../../resources -- go run .
```

No app port
Dapr port: 3500


## Scoreboard API

Scoreboard API Service is a Dapr service that saves volleyball game state,
and provides an API to retrieve game scores using Dapr topic event and service invocation handlers.
It listens to incoming game score update events on the `gamepubsub` pubsub `game` topic,
and any game score that is game point (25) or higher it will save to the statestore.
Specific game score may be found using this API when provided a game ID.

```
cd cmd/scoreboard

dapr run \
  --app-port 3002 \
  --app-id scoreboard \
  --app-protocol http \
  --dapr-http-port 3500 \
  --resources-path=../../resources -- go run .
```

App port 3002
Dapr port: 3500


## Game Service

Game Service is a Dapr service that provides an interface for the web UI to interact with the system.
It has a `scoreboard` endpoint that invokes service invocation on the `scoreboard` service to retrieve game score for a specific game ID to display on the web UI.

```
cd cmd/game-service

dapr run \
--app-id game-service \
--app-port 3001 \
--app-protocol http \
--dapr-http-port 3500 \
--resources-path ../../resources -- go run .
```

App Port: 3001
Dapr port: 3500

## Web UI

The Web UI displays volleyball game score information.

```
cd web-ui/

npm start
```

UI can be reached at: http://localhost:3000/