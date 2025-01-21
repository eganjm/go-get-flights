# Use the official Go image from Docker Hub
FROM golang:1.22-alpine AS builder

# Install dependencies for Go modules (optional for Alpine)
RUN apk add --no-cache git

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files and Go program
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Download Go module dependencies (if any)
RUN go mod tidy

# Build the Go program
RUN go build -o departures-service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/departures-service .

# Expose port 8080 for the web server
EXPOSE 8081

# Command to run the compiled Go binary
CMD ["./departures-service"]