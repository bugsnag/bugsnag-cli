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
	npm i && npm install -g npm-check && npm-check

.PHONY: go-lint
go-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	golangci-lint run

.PHONY: bump
bump:
ifneq ($(shell git diff --staged),)
	@git diff --staged
	@$(error You have uncommitted changes. Push or discard them to continue)
endif
	@./scripts/bump-version.sh $(VERSION)

.PHONY: test-fixtures
test-fixtures: features/base-fixtures/android features/base-fixtures/dart features/base-fixtures/rn0_69 features/base-fixtures/rn0_70 features/base-fixtures/rn0_72 features/base-fixtures/dsym/swift-package-manager

.PHONY: features/base-fixtures/android
features/base-fixtures/android:
	cd $@ && ./gradlew bundleRelease

.PHONY: features/base-fixtures/dart
features/base-fixtures/dart:
	cd $@ && $(FLUTTER_BIN) pub get
	cd $@ && $(FLUTTER_BIN) build apk  --suppress-analytics --split-debug-info=app-debug-info
	cd $@ && $(FLUTTER_BIN) build ios --no-codesign --suppress-analytics --no-tree-shake-icons --split-debug-info=app-debug-info

.PHONY: features/base-fixtures/rn0_69
features/base-fixtures/rn0_69: features/base-fixtures/rn0_69/android features/base-fixtures/rn0_69/ios

.PHONY: features/base-fixtures/rn0_70
features/base-fixtures/rn0_70: features/base-fixtures/rn0_70/android features/base-fixtures/rn0_70/ios

.PHONY: features/base-fixtures/rn0_72
features/base-fixtures/rn0_72: features/base-fixtures/rn0_72/android features/base-fixtures/rn0_72/ios

.PHONY: features/base-fixtures/rn0_69/android
features/base-fixtures/rn0_69/android:
	cd $@/../ && npm i
	cd $@ && ./gradlew bundleRelease

.PHONY: features/base-fixtures/rn0_69/ios
features/base-fixtures/rn0_69/ios:
	cd $@/../ && npm i && bundle install
	cd $@ && pod install
	cd $@ && xcodebuild -workspace rn0_69.xcworkspace -scheme rn0_69 -configuration Release -sdk iphoneos build

.PHONY: features/base-fixtures/rn0_70/android
features/base-fixtures/rn0_70/android:
	cd $@/../ && npm i
	cd $@ && ./gradlew bundleRelease

.PHONY: features/base-fixtures/rn0_70/ios
features/base-fixtures/rn0_70/ios:
	cd $@/../ && npm i && bundle install
	cd $@ && pod install
	cd $@ && xcodebuild -workspace rn0_70.xcworkspace -scheme rn0_70 -configuration Release -sdk iphoneos build

.PHONY: features/base-fixtures/rn0_72/android
features/base-fixtures/rn0_72/android:
	cd $@/../ && npm i
	cd $@ && ./gradlew bundleRelease

.PHONY: features/base-fixtures/rn0_72/ios
features/base-fixtures/rn0_72/ios:
	cd $@/../ && npm i && bundle install
	cd $@ && pod install
	cd $@ && xcodebuild -workspace rn0_72.xcworkspace -scheme rn0_72 -configuration Release -sdk iphoneos build

.PHONY: features/base-fixtures/dsym/swift-package-manager
features/base-fixtures/dsym/swift-package-manager:
	cd $@ && bundle install
	echo "--- Resolve Swift Package Dependencies"
	cd $@ && xcodebuild -allowProvisioningUpdates -scheme swift-package-manager -derivedDataPath DerivedData -resolvePackageDependencies
	echo "+++ Build Release iOS"
	cd $@ && xcodebuild -allowProvisioningUpdates -scheme swift-package-manager -configuration Release -destination generic/platform=iOS -derivedDataPath DerivedData -quiet build GCC_TREAT_WARNINGS_AS_ERRORS=YES
	#cd $@ && xcodebuild -scheme swift-package-manager -configuration Release -sdk iphoneos build -allowProvisioningUpdates
