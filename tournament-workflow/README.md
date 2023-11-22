# To run
setup venv and activate
source tournamentvenv/bin/activate
pip3 install -r requirements

dapr run --app-id orderapp --app-protocol grpc --dapr-grpc-port 50001 --resources-path components --placement-host-address localhost:50005 -- python3 flaskapp.py


# start workflow
#### todo fix and actually pass in the right data for payload
curl -X POST -H "Content-Type: application/json" -d '{"team_name": "Team A", "captain": "John", "players": ["Player1", "Player2"]}' http://127.0.0.1:5001/start-workflow

