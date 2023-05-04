# Dapr Volleyball Demo

## Volleyball Game Simulator
```
dapr run --app-id game-sim --app-protocol http --dapr-http-port 3500 --resources-path ../components -- go run .
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