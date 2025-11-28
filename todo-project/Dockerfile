# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go mod download
RUN go mod tidy
RUN go build -o todo-project main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-project .

ENV PORT=3000

CMD ["./todo-project"]
