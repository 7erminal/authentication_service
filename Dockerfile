# --- Stage 1: Build the Application ---
FROM golang:1.24-alpine AS builder

# Set the source workspace directory
WORKDIR /src

# Install critical system packages required for the Go compiler
RUN apk add --no-cache git ca-certificates

# Leverage Docker cache layers for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining project codebase into the container build frame
COPY . .

# Compile optimized, static binary dropping debugging symbols and layout info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/authentication_service .

# --- Stage 2: Production Distroless Alpine Runtime ---
FROM alpine:3.20

# Create application runtime environment directory context
WORKDIR /app

# Install standard localization/security assets, initialize non-root users/groups
RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S appgroup \
	&& adduser -S appuser -G appgroup \
	&& mkdir -p /logs /app/conf /app/swagger \
	&& chown -R appuser:appgroup /logs /app

# Copy artifact directly out of the build stage
COPY --from=builder /out/authentication_service /app/authentication_service

# Copy application structural context requirements into the active work directory
COPY conf /app/conf
COPY swagger /app/swagger

# Switch to the non-root execution context for enhanced safety boundaries
USER appuser

# Document that the app is listening on the targeted 5080 layout
EXPOSE 5080

# Execute relative to /app directory so Beego can automatically map ./conf/app.conf
CMD ["./authentication_service"]