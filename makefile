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

.PHONY: build
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
	GOOS=linux GOARCH=arm64 go build -ldflags '-s' -o bin/arm64-linux-bugsnag-cli main.go
	GOOS=linux GOARCH=amd64 go build -ldflags '-s' -o bin/x86_64-linux-bugsnag-cli main.go
	GOOS=linux GOARCH=386 go build -ldflags '-s' -o bin/i386-linux-bugsnag-cli main.go

.PHONY: build-macos
build-macos:
	[ -d bin ] || mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s' -o bin/x86_64-macos-bugsnag-cli main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags '-s' -o bin/arm64-macos-bugsnag-cli main.go

.PHONY: fmt
fmt:
	gofmt -w ./

.PHONY: unit-tests
unit-test:
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@v2.5.0
	go test -race -json -v ./test/... 2>&1 | tee /tmp/gotest.log | gotestfmt

.PHONY: npm-lint
npm-lint:
	cd js && npm i && npm install -g npm-check && npm-check

.PHONY: go-lint
go-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@vlatest
	golangci-lint run

.PHONY: bump
bump:
ifneq ($(shell git diff --staged),)
	@git diff --staged
	@$(error You have uncommitted changes. Push or discard them to continue)
endif
	@./scripts/bump-version.sh $(VERSION)

.PHONY: test-fixtures
test-fixtures: features/base-fixtures/android features/base-fixtures/dart features/base-fixtures/rn0_69 features/base-fixtures/rn0_70 features/base-fixtures/rn0_72

.PHONY: features/base-fixtures/android
features/base-fixtures/android:
	cd $@ && ./gradlew bundleRelease

.PHONY: features/base-fixtures/dart
features/base-fixtures/dart:
	cd $@ && $(FLUTTER_BIN) pub get
	cd $@ && $(FLUTTER_BIN) build apk  --suppress-analytics --split-debug-info=app-debug-info
	cd $@ && $(FLUTTER_BIN) build ios --no-codesign --suppress-analytics --no-tree-shake-icons --split-debug-info=app-debug-info

.PHONY: features/base-fixtures/dsym
features/base-fixtures/dsym:
	bundle install
	cd $@ && xcrun xcodebuild -allowProvisioningUpdates -scheme dSYM-Example -resolvePackageDependencies
	cd $@ && xcrun xcodebuild -allowProvisioningUpdates -scheme dSYM-Example -configuration Release -quiet build GCC_TREAT_WARNINGS_AS_ERRORS=YES

.PHONY: features/base-fixtures/dsym/archive
features/base-fixtures/dsym/archive:
	bundle install
	cd features/base-fixtures/dsym && bundle install
	cd features/base-fixtures/dsym && xcrun xcodebuild -allowProvisioningUpdates -scheme dSYM-Example -resolvePackageDependencies
	cd features/base-fixtures/dsym && xcrun xcodebuild -scheme dSYM-Example -configuration Release -allowProvisioningUpdates archive

.PHONY: features/base-fixtures/dsym/archive/export
features/base-fixtures/dsym/archive/export:
	bundle install
	cd features/base-fixtures/dsym && bundle install
	cd features/base-fixtures/dsym && xcrun xcodebuild -allowProvisioningUpdates -scheme dSYM-Example -resolvePackageDependencies
	cd features/base-fixtures/dsym && xcrun xcodebuild DEVELOPMENT_TEAM=7W9PZ27Y5F -scheme dSYM-Example -configuration Release -archivePath dsym.xcarchive -allowProvisioningUpdates archive
	cd features/base-fixtures/dsym && xcrun xcodebuild -exportArchive -archivePath dsym.xcarchive -exportPath output/ -exportOptionsPlist ../../react-native/fixtures/app/dynamic/ios/exportOptions.plist

.PHONY: features/base-fixtures/js-webpack4
features/base-fixtures/js-webpack4:
	cd $@ && npm i && npm run build

.PHONY: features/base-fixtures/js-webpack5
features/base-fixtures/js-webpack5:
	cd $@ && npm i && npm run build
