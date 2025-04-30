run:
	@go run cmd/main.go

test:
	@go test -v ./...

test-coverage:
	@go test -v ./... -cover
