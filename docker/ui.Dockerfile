# # Final Stage: Build UI and combine with microservices
# FROM node:16-alpine AS builder
# WORKDIR /app
# COPY web-ui/ .
# RUN npm install
# RUN npm run build

# # Expose port 80 for the application to listen on
# EXPOSE 80


# FROM node:16-alpine AS deploy
# COPY --from=builder /app/build /app/ui
# # COPY --from=builder /app/build /usr/share/nginx/html

# CMD "npm start"

# Final Stage: Build UI and combine with microservices
# FROM node:16-alpine AS builder
# WORKDIR /app
# COPY web-ui/ .
# ENV API_URL="http://game-service.default.svc.cluster.local:80"
# RUN npm install
# RUN npm run build

# # Expose port 80 for the application to listen on
# EXPOSE 80

# FROM node:16-alpine AS deploy
# WORKDIR /app
# COPY --from=builder /app/build ./ui
# RUN apk add --no-cache bash
# COPY web-ui/package*.json ./
# RUN npm install --only=production
# COPY web-ui .
# ENV API_URL="http://game-service.default.svc.cluster.local:80"
# CMD ["npm", "start"]




# FROM node:16-alpine
# WORKDIR /app
# COPY web-ui/ .
# ENV API_URL="http://game-service.default.svc.cluster.local:80"
# RUN npm install
# RUN npm run build

# Expose port 80 for the application to listen on

# FROM node:16-alpine AS deploy
# WORKDIR /app
# COPY --from=builder /app/build ./ui
# RUN apk add --no-cache bash
# COPY web-ui/package*.json ./
# RUN npm install --only=production
# COPY web-ui .
# ENV API_URL="http://game-service.default.svc.cluster.local:80"
# EXPOSE 80

# CMD ["npm", "start"]



# FROM node:16-alpine
# WORKDIR /app
# COPY web-ui/ .
# ENV API_URL="http://game-service.default.svc.cluster.local:80"
# RUN npm install --only=production
# EXPOSE 80
# CMD ["npm", "start"]


# Build stage
FROM node:16-alpine3.14 AS build
RUN apk add --no-cache git
WORKDIR /app
COPY web-ui/package*.json ./
RUN npm install
COPY web-ui/ ./
# RUN npm run build

# Final stage
# FROM alpine:3.14
# WORKDIR /app
# COPY --from=build /app/build/ ./ui/
ENV API_URL="http://game-service.default.svc.cluster.local:80"
EXPOSE 80
CMD ["npm", "start"]

# Build stage
# FROM node:16-alpine3.14 AS build
# RUN apk add --no-cache git
# WORKDIR /app
# COPY web-ui/package*.json ./
# RUN npm install 
# COPY web-ui/ ./
# RUN npm run build

# # Final stage
# FROM alpine:3.14
# WORKDIR /app
# COPY --from=build /app/build/ ./ui/
# ENV API_URL="http://game-service.default.svc.cluster.local:80"
# EXPOSE 80
# ENV PATH="/app/node_modules/.bin:${PATH}"

# CMD ["npm", "start"]

