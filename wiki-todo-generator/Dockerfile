FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
COPY main.go ./

RUN go mod download
RUN go build -o wiki-todo-generator main.go

FROM alpine:3.24

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/wiki-todo-generator .

CMD ["./wiki-todo-generator"]
