# Dapr Volleyball Demo

## Volleyball Game Simulator

### Run locally:

```
cd cmd/game-sim
go install ../../pkg/types.go 
go build .
dapr run --app-id game-sim --app-protocol http --dapr-http-port 3502 --resources-path ../../resources -- go run .
```

No app port
Dapr port: 3500


## Scoreboard API

```
cd cmd/scoreboard
go install ../../pkg/types.go 
go build .
dapr run \
  --app-port 6006 \
  --app-id scoreboard \
  --app-protocol http \
  --dapr-http-port 3501 \
  --resources-path=../../resources -- go run .
```

App port 6006 - 3002??
Dapr port: 3500


## Game Service
```
cd cmd/game-service
dapr run --app-id game-service --app-port 3001 --app-protocol http --dapr-http-port 3500 --resources-path ../../resources -- go run .
```

App Port: 3001
Dapr port: 3500

## Web UI
```
cd web-ui/
npm start
```

UI can be reached at: http://localhost:3000/