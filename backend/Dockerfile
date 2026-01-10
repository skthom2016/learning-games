# Backend Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

# Install air for hot reload during development
RUN go install github.com/air-verse/air@v1.61.0

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Air will be run via docker-compose command
CMD ["air", "-c", ".air.toml"]
