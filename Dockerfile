FROM ubuntu:22.04 AS builder
LABEL authors="gabriel"

RUN apt-get update && apt-get install -y \
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
