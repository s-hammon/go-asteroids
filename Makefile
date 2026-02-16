PROJECT_NAME := go-asteroids
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

build-bin: test clean
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/${PROJECT_NAME}

build-wasm: clean
	@GOOS=js GOARCH=wasm go build -o bin/${PROJECT_NAME}.wasm

test:
	@go vet ./...
	@go test -cover ./...

clean:
	@go mod tidy

gosec:
	@gosec -terse -exclude=G404 ./...

lint:
	@golangci-lint run --timeout=2m

ready: test lint gosec

test-packages:
	go test -json $$(go list ./... | grep -v -e /bin -e /cmd -e /internal/api/models) |\
		tparse --follow -sort=elapsed -trimpath=auto -all

.PHONY: clean test build-bin gosec lint ready
