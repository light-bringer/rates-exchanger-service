# Start from the official golang:alpine base image
FROM golang:1.22-bullseye AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    SERVICE_PATH=./cmd/api-service

# Move to working directory /build
WORKDIR /build

# Copy and download modules
COPY go.mod .
COPY go.sum .
RUN go mod vendor

# Copy the code into the container
COPY . .

RUN ls -la
RUN pwd

# Build the application
RUN go build -o rates-api $SERVICE_PATH

# Start a new stage from scratch
FROM gcr.io/distroless/static-debian11

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/rates-api /app/
COPY --from=builder /build/config.yaml /app/


# Command to run the executable
CMD ["/app/rates-api", "--config-file", "/app/config.yaml"]
