# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN go mod download
RUN go build -o ping-pong main.go

# Run stage
FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/ping-pong .

CMD ["./ping-pong"]
