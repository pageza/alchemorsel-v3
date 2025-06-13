.PHONY: run test

run:
	cd backend && go run ./cmd/api

test:
	cd backend && go test ./...
