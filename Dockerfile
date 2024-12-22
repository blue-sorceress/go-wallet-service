# Use the official Golang image as a base
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-wallet-service .

# Use a smaller base image for the final image
FROM ubuntu:latest

# Install PostgreSQL and other necessary tools
RUN apt-get update && \
    apt-get install -y postgresql postgresql-contrib && \
    apt-get clean

RUN apt-get update && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/*

# Copy the Go binary from the builder stage
COPY --from=builder /app/go-wallet-service /usr/local/bin/

# Create a directory for PostgreSQL data
RUN mkdir -p /var/lib/postgresql/data

# Start PostgreSQL and your Go application
CMD service postgresql start && \
    /usr/local/bin/go-wallet-service