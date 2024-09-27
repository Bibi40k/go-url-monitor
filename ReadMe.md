# Initialize the Go module
go mod init go-url-monitor
go mod tidy

# Run the application using Docker Compose
docker-compose up --build