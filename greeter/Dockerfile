# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN go build -o greeter main.go

# Run stage
FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/greeter .

CMD ["./greeter"]
