FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy all necessary directories first
COPY pkg/ /app/pkg/
COPY transaction-ledger/ /app/transaction-ledger/
COPY transaction-processment/ /app/transaction-processment/

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/transaction-processment/server ./transaction-processment/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/transaction-processment/server .

CMD ["./server"] 