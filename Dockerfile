# Stage 1: Build
FROM golang:1.24-alpine AS builder

# application path
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app

# Add certs for HTTPS (if needed)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
ENTRYPOINT ["./main"]