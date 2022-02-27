.PHONY: test
test: build
	./test.sh

.PHONY: build
build:
	go build -o bin/favicheck
