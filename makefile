all: build

.PHONY: build
build:
	go build -ldflags '-s'
