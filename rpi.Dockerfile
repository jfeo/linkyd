FROM golang:1.22.2-bullseye

RUN apt update && apt install -y gcc-arm-linux-gnueabihf

WORKDIR /linkyd

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc go build -ldflags="-s -w" -o build/linkyd-rpi .
