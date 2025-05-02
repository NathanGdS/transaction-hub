run-ledger:
	@go run transaction-ledger/cmd/main.go

run-processment:
	@go run transaction-processment/cmd/main.go

test:
	@go test -v ./...

test-coverage:
	@go test -v ./... -cover
