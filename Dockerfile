# Start from the official Golang image for building
FROM golang:1.21-alpine AS builder

WORKDIR /usr/app

# Install git for go mod and beego tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Beego application
RUN go build -o main .

# Use a minimal image for running
FROM alpine:latest

WORKDIR /usr/app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy the built binary from builder
COPY --from=builder /usr/app/main .

# Expose the port your Beego app listens on (default 8080)
EXPOSE 8080

# Run the application
CMD ["./main"]