# Dockerfile for Go project using Air
FROM golang:1.23-alpine

WORKDIR /app

# Install Air globally
RUN go install github.com/air-verse/air@latest

# Copy the Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Start Air for live reloading
CMD ["air"]