# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY main.go ./

# Download dependencies
RUN go mod download
RUN go mod tidy

# Build the application
RUN go build -o todo-backend main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/todo-backend .

# Expose port
EXPOSE 3000

# Run the application
CMD ["./todo-backend"]
