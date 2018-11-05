PRETTYTEST := $(shell command -v prettytest 2> /dev/null)

build:
	@go build -mod=vendor -o ./dist/cosmo ./cmd/cosmo

build-linux:
	@GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./dist/cosmo-linux-amd64 ./cmd/cosmo

test:
ifdef PRETTYTEST
	@go test -mod=vendor -cover ./... | prettytest
else
	@go test -mod=vendor -cover ./...
endif

clean:
	@rm -rf ./dist

.PHONY: build build-linux test clean
