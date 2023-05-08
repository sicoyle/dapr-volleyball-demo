.PHONY: docker-build

GO_SERVICES := game-service scoreboard game-sim

docker-build-ui:
	docker build -t sam-test-ui:latest --platform=linux/amd64 . -f docker/ui.Dockerfile

docker-build: docker-build-ui
	for service in $(GO_SERVICES); do \
		docker build -t $$service:latest --platform=linux/amd64 --build-arg SERVICE_NAME=$$service . -f docker/go.Dockerfile; \
	done

k8s-load-ui:
	kind load docker-image sam-test-ui:latest

k8s-load: k8s-load-ui
	for service in $(GO_SERVICES); do \
		kind load docker-image $$service:latest; \
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