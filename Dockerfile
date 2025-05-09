# Use the official Go image as the base image
FROM golang:1.24 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for dependency resolution
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum files are unchanged
RUN go mod download

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy the source code to the container
COPY . .

# Build the Go app
RUN go build -o currency_exchange ./cmd/api/

# Use a minimal base image for the final executable
FROM debian:bookworm-slim

# Install necessary dependencies, including PostgreSQL client if migrations or database seeding is needed
RUN apt-get update && apt-get install -y \
    ca-certificates \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables for the app
ENV APP_ENV=production \
    PORT=8080

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the prebuilt binary file from the builder
COPY --from=builder /app/currency_exchange .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./currency_exchange"]

