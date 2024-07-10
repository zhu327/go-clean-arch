FROM golang:1.21-bullseye

WORKDIR /code

COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

COPY . .