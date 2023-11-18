# Stage 1: Build stage
FROM golang:1.21.4-alpine AS build

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY server/ip_logger.go ip_logger.go

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ip_logger ip_logger

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/ip_logger .

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Set the entrypoint command
ENTRYPOINT ["/app/ip_logger"]
