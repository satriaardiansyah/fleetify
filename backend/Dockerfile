# ── Stage 1: Builder ──────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git (needed by some go modules)
RUN apk add --no-cache git

# Download dependencies first (layer cache)
COPY go.mod go.su[m] ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./app

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# ca-certificates for HTTPS calls, tzdata for Asia/Jakarta
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server .

# Copy config if it exists (ignored if missing via .dockerignore)
COPY --from=builder /app/config ./config

EXPOSE 8080

ENTRYPOINT ["./server"]
