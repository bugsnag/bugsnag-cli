all: build

.PHONY: build
build:
	[ -d build ] || mkdir -p bin
	go build -ldflags '-s'

.PHONY: build-all
build-all: build-windows build-linux build-mac

.PHONY: build-windows
build-windows:
	[ -d build ] || mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -ldflags '-s' -o bin/bugsnag-cli-x64.exe main.go
	GOOS=windows GOARCH=386 go build -ldflags '-s' -o bin/bugsnag-cli-386.exe main.go

.PHONY: build-linux
build-linux:
	[ -d build ] || mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags '-s' -o bin/bugsnag-cli-amd64-linux main.go
	GOOS=linux GOARCH=386 go build -ldflags '-s' -o bin/bugsnag-cli-386-linux main.go

.PHONY: build-mac
build-mac:
	[ -d build ] || mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s' -o bin/bugsnag-cli-amd64-darwin main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags '-s' -o bin/bugsnag-cli-arm64-darwin main.go