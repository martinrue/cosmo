build:
	@go build -mod=vendor -o ./dist/cosmo ./cmd/cosmo

build-linux:
	@GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./dist/cosmo-linux-amd64 ./cmd/cosmo

test:
	@go test -mod=vendor ./cosmo

clean:
	@rm -rf ./dist

.PHONY: build build-linux test clean
