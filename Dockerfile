FROM golang:1.22.2-alpine AS builder
LABEL authors="gabriel"

RUN apk update && apk add --no-cache \
    wkhtmltopdf \
    build-base \
    git \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .

RUN go build -o main ./main.go

RUN chmod +x main

RUN ls -l main
RUN whoami
RUN id

EXPOSE 8080

CMD ["./main"]

