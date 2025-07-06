# Dockerfile for the gateway service
# Stage 1 - Build the Go application
FROM golang:1.22.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY services/gateway/ ./
RUN go build -o gateway main.go

# Stage 2 - Create a minimal image with the built application
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gateway .
EXPOSE 80
CMD ["./gateway"]

