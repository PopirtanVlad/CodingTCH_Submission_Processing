FROM golang:1.18.2

RUN apt-get update && apt-get install -y \
    python3 \
    openjdk-11-jdk-headless

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

CMD ["./main"]