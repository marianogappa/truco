.PHONY: test build run release lint

test:
    go test -v ./...

build:
    go build -o truco ./...

run:
    ./truco

release:
    goreleaser --snapshot --rm-dist

lint:
    golangci-lint run