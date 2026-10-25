# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN go mod download
RUN go build -o log-writer main.go

# Run stage
FROM alpine:3.24

# Install timezone data
RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/log-writer .

CMD ["./log-writer"]
