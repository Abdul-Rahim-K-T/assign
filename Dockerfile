# Step 1: Use the official Golang image to build the Go app
FROM golang:1.22.1-alpine as builder

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code into the container
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

# Step 2: Create a minimal final image with only the built Go binary
FROM alpine:latest

# Set the working directory for the container
WORKDIR /root/

# Install necessary dependencies for running the Go app
RUN apk add --no-cache libc6-compat

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file (optional, depending on how you want to pass environment variables)
COPY .env .env

# Expose port 8080 for the application
EXPOSE 8080

# Command to run the Go binary when the container starts
CMD ["./main"]
