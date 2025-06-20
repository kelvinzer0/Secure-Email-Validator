# Build stage
FROM golang:1.21-alpine AS builder

# Install dig for DNSSEC validation
RUN apk add --no-cache bind-tools git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Secure Email Validator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o secure-email-validator .

# Final stage
FROM alpine:latest

# Install dig and ca-certificates for security validation
RUN apk --no-cache add bind-tools ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/secure-email-validator .

# Expose port for API server mode
EXPOSE 8080

# Health check for service monitoring
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ./secure-email-validator -server -port 8080 & sleep 1 && \
      wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command shows help
CMD ["./secure-email-validator", "-help"]
