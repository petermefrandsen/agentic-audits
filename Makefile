.PHONY: build test coverage coverage-report clean dist

# Ensure /usr/local/go/bin is in PATH (common location)
export PATH := $(PATH):/usr/local/go/bin

# Try to find go binary, fallback to 'go'
GO := $(shell which go 2>/dev/null || echo go)
BINARY_NAME=agentic-audits



build:
	$(GO) build -o $(BINARY_NAME) src/*.go

test:
	$(GO) test -v -coverprofile=coverage.out ./src/...
	$(GO) tool cover -func=coverage.out


coverage: test
	$(GO) tool cover -html=coverage.out

coverage-report: test
	@chmod +x src/generate_coverage_report.sh
	@./src/generate_coverage_report.sh


clean:
	rm -f $(BINARY_NAME) coverage.out

dist: build
	mkdir -p dist
	cp $(BINARY_NAME) dist/
	# Note: In a real scenario, you'd build for multiple platforms
	# GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64
