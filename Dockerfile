# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ ./

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -a -installsuffix cgo \
    -o web-streaming .

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create app user (security best practice)
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser

# Copy binary from builder
COPY --from=builder /build/web-streaming /bin/web-streaming

# Set working directory
WORKDIR /app

# Set environment
ENV GIN_MODE=release \
    PORT=1010

EXPOSE 1010

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:1010/health || exit 1

ENTRYPOINT ["/bin/web-streaming"]
