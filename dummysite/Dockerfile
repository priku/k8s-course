# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN CGO_ENABLED=0 go build -o dummysite-controller main.go

# Run stage
FROM alpine:3.24

WORKDIR /app

COPY --from=builder /app/dummysite-controller .

CMD ["./dummysite-controller"]
