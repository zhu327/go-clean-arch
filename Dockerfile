# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server ./cmd/api/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app -h /app app

WORKDIR /app
COPY --from=builder --chown=app:app /server /app/server

USER app
EXPOSE 8000
# The runtime PORT may override this default; Compose publishes HOST_PORT to PORT.

ENTRYPOINT ["/app/server"]
