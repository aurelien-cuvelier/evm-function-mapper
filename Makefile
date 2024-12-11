BIN_NAME=evmfm
.DEFAULT_GOAL := build

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BIN_NAME}-linux ./cmd/main.go
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BIN_NAME}-mac ./cmd/main.go
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BIN_NAME}-windows ./cmd/main.go