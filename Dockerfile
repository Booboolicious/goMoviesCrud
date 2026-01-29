# Build stage
FROM golang:1.25.0-bookworm AS builder

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Use CGO_ENABLED=1 because SQLite requires C
RUN CGO_ENABLED=1 go build -o movies-api ./Library

# Run stage
FROM debian:bookworm-slim

# Install sqlite3 for management (optional) and ca-certificates
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/movies-api .

# Copy the frontend directory
COPY --from=builder /app/frontend/ ./frontend/

# Expose the application port
EXPOSE 8000

# Run the binary
CMD ["./movies-api"]
