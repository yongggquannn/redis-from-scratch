FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy Go modules and dependencies
COPY go.mod ./
# go.sum may not exist initially; mod download will create it
RUN go mod download

# Copy source
COPY . .

# Build application
RUN go build -o server .

# Runtime image
FROM alpine:latest AS runtime
WORKDIR /app

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

# Copy binary
COPY --from=builder /app/server /usr/local/bin/server

# Make workdir writable for AOF file
RUN chown -R app:app /app
USER app

EXPOSE 6379

ENTRYPOINT ["/usr/local/bin/server"]
