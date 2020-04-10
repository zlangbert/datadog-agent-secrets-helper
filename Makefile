BINARY        ?= datadog-agent-secrets-helper
SOURCES        = $(shell find . -name '*.go')
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -w -s

.PHONY: build
build: build/$(BINARY)

build/$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

.PHONY: format
format:
	go fmt ./...