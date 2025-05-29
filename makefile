.PHONY: run


run:
	go run cmd/server/main.go

test_repository:
	go test ./internal/repository

coverage_repository:
	go test ./internal/repository -coverprofile=$@.out
	go tool cover -html=$@.out