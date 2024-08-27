# Setup go backend
FROM golang:1.23 AS build_env

WORKDIR /opt/server

COPY server/ ./
RUN go mod tidy && go mod download
WORKDIR /opt/server/cmd/gobloks
RUN CGO_ENABLED=0 GOOS=linux go build -o /opt/docker-gobloks-server

# Setup Vue frontend environment
FROM node:20.12-bookworm AS node_build_env

WORKDIR /tmp/frontend
COPY frontend/ .

RUN npm ci && npm run build

# Package up the final image
FROM alpine:3.20

# Install necessary packages
RUN apk add --no-cache nginx supervisor

WORKDIR /opt

# Copy frontend build files
COPY --from=node_build_env /tmp/frontend/dist /usr/share/nginx/html

# Copy Go application
COPY --from=build_env /opt/docker-gobloks-server ./docker-gobloks-server

COPY nginx.conf /etc/nginx/http.d/default.conf

COPY init.conf /etc/supervisord.conf

EXPOSE 7777 8888

CMD ["supervisord", "-c", "/etc/supervisord.conf"]
