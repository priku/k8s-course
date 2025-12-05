# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go mod download
RUN go mod tidy
RUN go build -o log-writer main.go

# Run stage
FROM alpine:latest

# Install timezone data
RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/log-writer .

CMD ["./log-writer"]
