FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go mod download
RUN go mod tidy
RUN go build -o wiki-todo-generator main.go

FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/wiki-todo-generator .

CMD ["./wiki-todo-generator"]
