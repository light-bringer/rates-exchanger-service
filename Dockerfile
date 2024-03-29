# Start from the official golang:alpine base image
FROM golang:alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    SERVICE_PATH=./cmd/api-service

# Move to working directory /build
WORKDIR /build

# Copy and download modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o rates-api $SERVICE_PATH

# Start a new stage from scratch
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/rates-api /app/

# Command to run the executable
CMD ["/app/rates-api"]
