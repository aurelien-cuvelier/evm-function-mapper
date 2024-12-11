BIN_NAME=evmfm
.DEFAULT_GOAL := build-all

build-all: build-linux build-mac build-windows

build-linux:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BIN_NAME}-linux ./cmd/main.go

build-mac:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BIN_NAME}-mac ./cmd/main.go

build-windows:
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BIN_NAME}-windows.exe ./cmd/main.go