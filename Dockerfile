FROM golang:1.22 AS builder
LABEL authors="gabriel"

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    tzdata \
    wkhtmltopdf \
    build-essential \
    git \
    && apt-get clean

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .

RUN go build -o main ./main.go

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]
