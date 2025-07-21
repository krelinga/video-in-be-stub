# Build stage
FROM golang:1.24.3-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o video-in-be-stub

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/video-in-be-stub .

# Change ownership to non-root user
RUN chown appuser:appgroup video-in-be-stub

# Switch to non-root user
USER appuser

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["./video-in-be-stub"]