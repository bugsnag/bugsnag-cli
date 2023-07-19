# Set PLATFORM var, if needed
ifeq ($(PLATFORM),)
ifeq ($(OS),Windows_NT)
PLATFORM=windows
else
UNAME := $(shell uname -s)
ifeq ($(UNAME),Linux)
PLATFORM=linux
endif
ifeq ($(UNAME),Darwin)
PLATFORM=macos
endif
endif
endif

PACKAGE_VERSION :=$(shell cat package.json | jq -r .version)

build: build-$(PLATFORM) # Build for PLATFORM or the host OS

.PHONY: build-all
build-all: build-windows build-linux build-macos

.PHONY: build-windows
build-windows:
	[ -d bin ] || mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -ldflags '-s -X main.package_version=$(PACKAGE_VERSION)' -o bin/x86_64-windows-bugsnag-cli.exe main.go
	GOOS=windows GOARCH=386 go build -ldflags '-s -X main.package_version=$(PACKAGE_VERSION)' -o bin/i386-windows-bugsnag-cli.exe main.go

.PHONY: build-linux
build-linux:
	[ -d bin ] || mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -X main.package_version=$(PACKAGE_VERSION)' -o bin/x86_64-linux-bugsnag-cli main.go
	GOOS=linux GOARCH=386 go build -ldflags '-s -X main.package_version=$(PACKAGE_VERSION)' -o bin/i386-linux-bugsnag-cli main.go

.PHONY: build-macos
build-macos:
	[ -d bin ] || mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -X main.package_version=$(PACKAGE_VERSION)' -o bin/x86_64-macos-bugsnag-cli main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags '-s -X main.package_version=$(PACKAGE_VERSION)' -o bin/arm64-macos-bugsnag-cli main.go

.PHONY: unit-test
unit-test:
	go test -json -v ./test/... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: fmt
fmt:
	gofmt -w ./
