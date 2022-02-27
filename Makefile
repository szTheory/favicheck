.PHONY: build
build:
	go build -o bin/favicheck

.PHONY: test
test: build
	./test.sh

lint: build
	golangci-lint run