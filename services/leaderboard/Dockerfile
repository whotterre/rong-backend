# Leaderboard multistage Dockerfile 
# Stage 1: Building the binary
FROM golang:1.22.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Awesome! 
# So the -ldflags reduces the size of the binary by
# stripping debugging information and symbol tables.
RUN CGO_ENABLED=0 go build -o main -ldflags "-s -w" .

# Stage 2: Creating the final image
FROM alpine:latest

# Similar to sudo apt update && sudo apt install ca-certificates
# On Ubuntu
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3001

CMD ["./main"]