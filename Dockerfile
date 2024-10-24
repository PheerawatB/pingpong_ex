# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code to the container
COPY . .

# Build the Go application
RUN go build -o pingpong-ex ./main.go

# Stage 2: Create a smaller container for running the app
FROM alpine:latest

# Set working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/pingpong-ex .

# Expose the ports for both services
EXPOSE 8888
EXPOSE 8889

# Command to run the application (which runs both services)
CMD ["./pingpong-ex"]
