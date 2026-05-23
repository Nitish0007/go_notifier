# Stage 1: Builder
FROM golang:1.25-alpine AS builder

ENV ENV=production
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o server ./cmd/server
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o workers ./cmd/workers

# Stage 2: API Runtime
FROM alpine:3.23 AS api

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

# Non-root user
RUN addgroup -S myuser && adduser -S myuser -G myuser
USER myuser

EXPOSE 8080

ENTRYPOINT ["./server"]

# Stage 3: Workers Runtime
FROM alpine:3.23 AS workers

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/workers .
COPY --from=builder /app/configs ./configs

# Non-root user
RUN addgroup -S myuser && adduser -S myuser -G myuser
USER myuser

ENTRYPOINT ["./workers"]


# command to build api image -> docker build --target api -t notifier-api .
# command to build workers image -> docker build --target workers -t notifier-workers .