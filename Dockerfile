# Build stage
FROM golang:1.25 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# # Install necessary tools and download dependencies
RUN apt-get update && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/* && \
    go install github.com/air-verse/air@latest && \
    go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the binaries
RUN GOOS=linux go build -o bin/go-backend-api-server cmd/api/main.go

# Run stage
FROM golang:1.25

# Set the working directory
WORKDIR /app

# Copy necessary files from the build stage
COPY --from=builder /app/bin /app/bin
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY .env .env

# Expose necessary ports
EXPOSE 8080

# Run the main API server with air
CMD ["air", "-c", ".air.toml"]
