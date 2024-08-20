# Setup go backend
FROM golang:1.23 AS build_env

WORKDIR /opt/server
COPY server/go.mod ./
RUN go mod download

COPY server/ ./
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

RUN echo -e \
"server { \n"\
" listen 7777 default_server; \n"\
" listen [::]:7777 default_server; \n"\
" root /usr/share/nginx/html; \n"\
" index index.html index.htm; \n"\
" location / { \n"\
"  try_files \$uri \$uri/ /index.html; \n"\
" } \n"\
"}"\
> /etc/nginx/http.d/default.conf

RUN echo -e \
"[supervisord] \n"\
"nodaemon=true \n"\
"[supervisorctl] \n"\
"[program:nginx] \n"\
"command=/usr/sbin/nginx -g 'daemon off;' \n"\
"[program:gobloks-server] \n"\
"command=/opt/docker-gobloks-server \n"\
> /etc/supervisord.conf

EXPOSE 7777 8888

CMD ["supervisord", "-c", "/etc/supervisord.conf"]
