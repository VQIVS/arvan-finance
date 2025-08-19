.PHONY: build run-api run-consumer run dev test clean deps fmt vet

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/consumer cmd/consumer/main.go

test:
	./tests/run_tests.sh

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

run-api:
	go run cmd/api/main.go

run-consumer:
	go run cmd/consumer/main.go

run-dev:
	$(MAKE) build && ($(MAKE) run-api & $(MAKE) run-consumer)