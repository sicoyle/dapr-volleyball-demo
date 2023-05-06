# Dapr Volleyball Demo

## Volleyball Game Simulator
```
dapr run --app-id game-sim --app-protocol http --dapr-http-port 3502 --resources-path ../components -- go run .
```

## Scoreboard API
```
dapr run \
  --app-port 6006 \
  --app-id scoreboard \
  --app-protocol http \
  --dapr-http-port 3501 \
  --resources-path=../components -- go run .
```

## Game Service
```
dapr run --app-id gameservice --app-protocol http --dapr-http-port 3500 --resources-path ../components -- go run .
```

## Web UI
```
cd web-ui/
npm start
```

UI can be reached at: http://localhost:3000/