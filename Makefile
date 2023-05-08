.PHONY: docker-build \
k8s-load-ui \
k8s-load \
k8s-deploy-dapr-resources \
k8s-clean-dapr-resources \
k8s-clean

GO_SERVICES := game-service scoreboard game-sim
DOCKER_TAG := latest

# Check for the DOCKER_REGISTRY environment variable
ifndef DOCKER_REGISTRY
$(error DOCKER_REGISTRY is undefined)
endif

docker-build:
	docker build -t $(DOCKER_REGISTRY)/web-ui:$(DOCKER_TAG) --platform=linux/amd64 . -f docker/ui.Dockerfile
	for service in $(GO_SERVICES); do \
		docker build -t $(DOCKER_REGISTRY)/$$service:$(DOCKER_TAG) --platform=linux/amd64 --build-arg SERVICE_NAME=$$service . -f docker/go.Dockerfile; \
	done

docker-push: docker-build
	docker push $(DOCKER_REGISTRY)/web-ui:$(DOCKER_TAG)
	for service in $(GO_SERVICES); do \
		docker push $(DOCKER_REGISTRY)/$$service:$(DOCKER_TAG); \
	done

k8s-load-ui:
	kind load docker-image web-ui:$(DOCKER_TAG)

k8s-load: k8s-load-ui
	for service in $(GO_SERVICES); do \
		kind load docker-image $$service:$(DOCKER_TAG); \
	done

k8s-deploy-dapr-resources:
	kubectl apply -f resources/k8s/pubsub-k8s.yaml
	kubectl apply -f resources/k8s/statestore-k8s.yaml

k8s-clean-dapr-resources:
	kubectl delete -f resources/k8s/pubsub.yaml
	kubectl delete -f resources/k8s/statestore.yaml

k8s-deploy: k8s-deploy-dapr-resources
	kubectl apply -f ./deploy/web-ui.yaml
	for service in $(GO_SERVICES); do \
		kubectl apply -f ./deploy/$$service.yaml; \
	done 

k8s-clean: k8s-clean-dapr-resources
	kubectl delete -f ./deploy/web-ui.yaml
	for service in $(GO_SERVICES); do \
		kubectl delete -f ./deploy/$$service.yaml; \
	done 