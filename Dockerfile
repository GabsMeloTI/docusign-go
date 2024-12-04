FROM golang:1.22.2-alpine AS builder
LABEL authors="gabriel"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go

RUN chmod +x main

RUN ls -l main
RUN whoami
RUN id

EXPOSE 8080

CMD ["./main"]
