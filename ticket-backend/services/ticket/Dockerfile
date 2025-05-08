# Stage 1: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git curl bash

# Install golang-migrate
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then ARCH=linux-amd64; \
    elif [ "$ARCH" = "aarch64" ]; then ARCH=linux-arm64; \
    else echo "Unsupported architecture: $ARCH" && exit 1; fi && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.${ARCH}.tar.gz \
    | tar -xz && mv migrate /usr/bin/migrate && chmod +x /usr/bin/migrate

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server

# Stage 2: Runtime image with migrate + Go binary
FROM alpine

WORKDIR /app

RUN apk add --no-cache bash curl

# Install golang-migrate
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then ARCH=linux-amd64; \
    elif [ "$ARCH" = "aarch64" ]; then ARCH=linux-arm64; \
    else echo "Unsupported architecture: $ARCH" && exit 1; fi && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.${ARCH}.tar.gz \
    | tar -xz && mv migrate /usr/bin/migrate && chmod +x /usr/bin/migrate

COPY --from=builder /app/server .

# Copy migration files
COPY db/migrations ./db/migrations

EXPOSE 8080 9090
