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

FLUTTER_BIN?=flutter

build: build-$(PLATFORM) # Build for PLATFORM or the host OS

.PHONY: build-all
build-all: build-windows build-linux build-macos

.PHONY: build-windows
build-windows:
	[ -d bin ] || mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -ldflags '-s' -o bin/x86_64-windows-bugsnag-cli.exe main.go
	GOOS=windows GOARCH=386 go build -ldflags '-s' -o bin/i386-windows-bugsnag-cli.exe main.go

.PHONY: build-linux
build-linux:
	[ -d bin ] || mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags '-s' -o bin/x86_64-linux-bugsnag-cli main.go
	GOOS=linux GOARCH=386 go build -ldflags '-s' -o bin/i386-linux-bugsnag-cli main.go

.PHONY: build-macos
build-macos:
	[ -d bin ] || mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s' -o bin/x86_64-macos-bugsnag-cli main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags '-s' -o bin/arm64-macos-bugsnag-cli main.go

.PHONY: unit-test
unit-test:
	go test -json -v ./test/... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: fmt
fmt:
	gofmt -w ./

bump:
ifneq ($(shell git diff --staged),)
	@git diff --staged
	@$(error You have uncommitted changes. Push or discard them to continue)
endif
	@./scripts/bump-version.sh $(VERSION)


test-fixtures: android-test-fixture dart-test-fixture rn0-69-test-fixture rn0-70-test-fixture rn0-72-test-fixture

android-test-fixture:
	cd features/base-fixtures/android && ./gradlew bundleRelease

dart-test-fixture:
	cd features/base-fixtures/dart && $(FLUTTER_BIN) pub get && $(FLUTTER_BIN) build apk  --suppress-analytics --split-debug-info=app-debug-info && $(FLUTTER_BIN) build ios --no-codesign --suppress-analytics --no-tree-shake-icons --split-debug-info=app-debug-info

rn0-69-test-fixture:
	cd features/base-fixtures/rn0_69 && npm i && cd android && ./gradlew bundleRelease

rn0-70-test-fixture:
	cd features/base-fixtures/rn0_70 && npm i && cd android && ./gradlew bundleRelease

rn0-72-test-fixture:
	cd features/base-fixtures/rn0_72 && npm i && cd android && ./gradlew bundleRelease
