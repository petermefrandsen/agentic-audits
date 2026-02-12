.PHONY: build test coverage clean dist

BINARY_NAME=agentic-audits

build:
	go build -o $(BINARY_NAME) src/*.go

test:
	go test -v -coverprofile=coverage.out ./src/...
	go tool cover -func=coverage.out

coverage: test
	go tool cover -html=coverage.out

clean:
	rm -f $(BINARY_NAME) coverage.out

dist: build
	mkdir -p dist
	cp $(BINARY_NAME) dist/
	# Note: In a real scenario, you'd build for multiple platforms
	# GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64
