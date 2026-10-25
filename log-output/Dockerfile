# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN go mod download
RUN go build -o log-output main.go

# Run stage
FROM alpine:3.24

# Install timezone data
RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/log-output .

CMD ["./log-output"]
