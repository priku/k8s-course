# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN go build -o wiki-todo-generator main.go

# Run stage
FROM alpine:3.24

# CA certificates for the HTTPS request to Wikipedia
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/wiki-todo-generator .

CMD ["./wiki-todo-generator"]
