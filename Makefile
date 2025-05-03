run-ledger:
	@go run transaction-ledger/cmd/main.go

run-processment:
	@go run transaction-processment/cmd/main.go

test:
	@go test -v ./...

test-coverage:
	@go test -v ./... -cover

loader-windows:
	@powershell -ExecutionPolicy Bypass -File loader.ps1

loader-linux:
	@./loader.sh

