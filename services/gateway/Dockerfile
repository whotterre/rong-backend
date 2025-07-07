# Dockerfile for the gateway service
# Stage 1 - Build the Go binary
FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o gateway main.go

# Stage 2 - Create a minimal image with the built application
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gateway .
EXPOSE 80
CMD ["./gateway"]

