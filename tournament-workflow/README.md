# To run
setup venv and activate
pip3 install -r requirements

dapr run --app-id orderapp --app-protocol grpc --dapr-grpc-port 50001 --resources-path components --placement-host-address localhost:50005 -- python3 flaskapp.py


