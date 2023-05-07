.PHONY: docker-build

GO_SERVICES := game-service scoreboard game-sim

docker-build:
	for service in $(GO_SERVICES); do \
		docker build -t $$service:latest --build-arg SERVICE_NAME=$$service . -f docker/go.Dockerfile; \
	done

k8s-load:
	for service in $(GO_SERVICES); do \
		kind load docker-image $$service:latest; \
	done

k8s-deploy-dapr-resources:
	kubectl apply -f resources/pubsub-k8s.yaml
	kubectl apply -f resources/statestore-k8s.yaml

k8s-clean-dapr-resources:
	kubectl delete -f resources/pubsub.yaml
	kubectl delete -f resources/statestore.yaml

k8s-deploy: k8s-deploy-dapr-resources
	for service in $(GO_SERVICES); do \
		kubectl apply -f ./deploy/$$service.yaml; \
	done 

k8s-clean: k8s-clean-dapr-resources
	for service in $(GO_SERVICES); do \
		kubectl delete -f ./deploy/$$service.yaml; \
	done 