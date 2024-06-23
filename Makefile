.PHONY: test build run release lint

test:
	go test -v ./...

build:
	go build -o truco ./...

run:
	./truco

release:
	rm -rf dist && goreleaser

lint:
	golangci-lint run