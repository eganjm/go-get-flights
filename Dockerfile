# Use the official Go image from Docker Hub
FROM golang:alpine

# Install dependencies for Go modules (optional for Alpine)
RUN apk add --no-cache git

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files and Go program
COPY go.mod .
COPY get-flights.go .

# Download Go module dependencies (if any)
RUN go mod tidy

# Build the Go program
RUN go build -o get-flights .

# Expose port 8080 for the web server
EXPOSE 8080

# Command to run the compiled Go binary
CMD ["./get-flights"]

