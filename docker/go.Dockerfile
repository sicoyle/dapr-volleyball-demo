# build stage
FROM golang:1.17.2-alpine3.14 AS builder
ARG SERVICE_NAME
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./pkg ./pkg
COPY ./cmd/${SERVICE_NAME}/ ./app/${SERVICE_NAME}/

RUN go build -o ${SERVICE_NAME} ./app/${SERVICE_NAME}/main.go

# Final image
FROM alpine:latest

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Copy the binary from the builder image
COPY --from=builder ./app/${SERVICE_NAME} ./${SERVICE_NAME}

# Make the binary executable
RUN chmod +x ./${SERVICE_NAME}

EXPOSE 80

# Start the service
CMD "./${SERVICE_NAME}"