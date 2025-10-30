# Stage 1: Build
FROM golang:1.24-alpine AS builder

# application path
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# build server
RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

# build workers
RUN GOOS=linux GOARCH=amd64 go build -o workers ./cmd/workers/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app

# Add certs for HTTPS (if needed)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/workers .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
# ENTRYPOINT ["./main"]