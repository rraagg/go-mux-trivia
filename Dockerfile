# syntax=docker/dockerfile:1
FROM golang:1.22.0-alpine3.14

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
