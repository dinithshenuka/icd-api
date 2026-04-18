# --- STAGE 1: Build ---
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the tiny production binary
# CGO_ENABLED=0 makes it a static binary (works perfectly with alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o icd-api ./cmd/api

# --- STAGE 2: Run ---
FROM alpine:latest

# Basic maintenance
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/icd-api .
# Copy the database (for a self-contained ready-to-run image)
COPY --from=builder /app/database/icd11.db ./database/icd11.db
# Copy the OpenAPI spec for the docs
COPY --from=builder /app/openapi/openapi.yaml ./openapi/openapi.yaml

# Set environment defaults
ENV PORT=8080
ENV GIN_MODE=release
ENV DB_PATH=database/icd11.db

EXPOSE 8080

CMD ["./icd-api"]
