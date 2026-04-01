
# Project Metadata
BINARY=ghost
BUILD_DIR=bin
LDFLAGS=-ldflags="-s -w"

# The Cross-Platform Duo
PLATFORMS=windows/amd64 linux/amd64 linux/arm64 linux/386 linux/arm/v7


.PHONY: all deps build-all clean

all: deps build-all

deps:
	@echo "👻 Sniffing out dependencies..."
	go mod tidy

build-all:
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d/ -f1); \
		ARCH=$$(echo $$platform | cut -d/ -f2); \
		EXTENSION=""; \
		if [ "$$OS" = "windows" ]; then EXTENSION=".exe"; fi; \
		OUT=$(BUILD_DIR)/$(BINARY)-$$OS-$$ARCH$$EXTENSION; \
		echo "🏗️  Building Ghost for $$platform..."; \
		GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=0 go build $(LDFLAGS) -o $$OUT .; \
	done

clean:
	rm -rf $(BUILD_DIR)
	rm -f helpers_*.go
