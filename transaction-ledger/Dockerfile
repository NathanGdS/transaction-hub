FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy pkg directory first
COPY pkg/ /app/pkg/

# Copy the transaction-ledger application
COPY transaction-ledger/ /app/transaction-ledger/

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/transaction-ledger/server ./transaction-ledger/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/transaction-ledger/server .

EXPOSE 8080

CMD ["./server"] 