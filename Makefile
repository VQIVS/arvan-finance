
.PHONY: build run run-dev run-api run-consumer test clean swagger docker-build docker-run docker-stop lint lint-fix lint-detailed check install-tools security

build:
	go build -o ./bin/api ./cmd/api
	go build -o ./bin/consumer ./cmd/consumer

test:
	go test -v ./...

clean:
	rm -rf ./bin

swagger:
	swag init -g ./cmd/api/main.go -o ./docs

fmt:
	go fmt ./...

vet:
	go vet ./...

install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

lint:
	golangci-lint run

lint-detailed:
	golangci-lint run --enable-all --disable=gochecknoglobals,gochecknoinits,godot,gomnd,gomodguard,goerr113,wrapcheck,exhaustruct,ireturn,varnamelen,nosnakecase

lint-fix:
	golangci-lint run --fix
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -local sms-dispatcher -w .; \
	else \
		echo "goimports not found. Run 'make install-tools' to install it."; \
	fi

security:
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Run 'make install-tools' to install it."; \
	fi

check: fmt vet lint test
	@echo "All checks passed!"

deps:
	go mod download

update-deps:
	go get -u ./...
	go mod tidy

run-api:
	go run ./cmd/api/main.go

run-consumer:
	go run ./cmd/consumer/main.go

run-dev:
	$(MAKE) build && $(MAKE) swagger && ($(MAKE) run-api & $(MAKE) run-consumer)

docker-build:
	docker build -t finance-service -f build/Dockerfile .

docker-logs:
	docker-compose logs -fn100 

docker-clean:
	docker system prune -f
