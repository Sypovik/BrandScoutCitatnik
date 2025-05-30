.PHONY: run test_repository coverage_repository test_handler coverage_handler


run:
	go run cmd/server/main.go

test_repository:
	go test ./internal/repository

coverage_repository:
	go test ./internal/repository -coverprofile=$@.out
	go tool cover -html=$@.out


test_handler:
	go test ./internal/handlers

coverage_handler:
	go test ./internal/handlers -coverprofile=$@.out
	go tool cover -html=$@.out


test_router:
	go test ./internal/router

coverage_router:
	go test ./internal/router -coverprofile=$@.out
	go tool cover -html=$@.out

test_all:
	go test -v ./internal/handlers/ ./internal/repository/ ./internal/router/
	

coverage_all:
	go test -v -coverprofile=$@.out ./internal/handlers/ ./internal/repository/ ./internal/router/
	go tool cover -html=$@.out