# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN go build -o broadcaster main.go

# Run stage
FROM alpine:3.24

# CA certificates for webhooks over HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/broadcaster .

CMD ["./broadcaster"]
