# ------ Build ------
FROM golang:1.25 AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN GOOS=linux CGO_ENABLED=1 go build -ldflags="-s -w" -o linkyd .

# ------ Run ------
FROM debian:bookworm-slim

WORKDIR /linkyd

COPY --from=builder /build/linkyd /linkyd/linkyd

ENTRYPOINT ["/linkyd/linkyd"]

