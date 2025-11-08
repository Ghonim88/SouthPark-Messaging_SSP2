#Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates 
WORKDIR /app

COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080
# Run the application
CMD ["./main"]