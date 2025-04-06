# Stage 1: Build Goravel (Go Application)
FROM golang:alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0  \
    GOARCH="amd64" \
    GOOS=linux

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the application
RUN go build --ldflags "-extldflags -static" -o main .

# Stage 2: Final Image
FROM alpine:latest

USER root

# clean up
RUN rm -rf /var/cache/apk/*

# Set working directory
WORKDIR /www

# Copy application binary and required files
COPY --from=builder /build/main /www/
COPY --from=builder /build/database/ /www/database/
COPY --from=builder /build/public/ /www/public/
COPY --from=builder /build/storage/ /www/storage/   

USER root
# Expose necessary ports
EXPOSE 3000 50051