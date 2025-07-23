# Use the official Golang image to build the application
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
# Copy the Go module files and the vendor directory
COPY go.mod go.sum ./
COPY vendor ./vendor

# Copy the rest of the application source code
COPY . .

# Build the application for the Lambda environment using the vendored dependencies
# The handler is in cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /main -tags lambda.norpc cmd/server/main.go

# Use a lightweight base image for the final container
FROM --platform=linux/amd64 alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /main /main

# Set the command to run when the container starts
CMD ["/main"]