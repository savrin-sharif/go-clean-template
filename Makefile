APP_NAME ?= go-clean-template
BIN ?= bin/server

.PHONY: tidy fmt test lint run build docker-build docker-up clean

tidy:
	go mod tidy

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.tmp/*')

test:
	go test ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/server

build:
	mkdir -p bin
	go build -o $(BIN) ./cmd/server

docker-build:
	docker build -t $(APP_NAME) .

docker-up:
	docker compose up --build

clean:
	rm -rf bin coverage.txt coverage.html
