FROM node:16-alpine3.14
WORKDIR /app
COPY web-ui/package*.json ./
RUN npm install
COPY web-ui/ ./
RUN npm run build --production
RUN npm install -g serve
EXPOSE 80
CMD serve -s build