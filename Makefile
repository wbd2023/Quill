SHELL := /usr/bin/env bash

QUILL_CMD = go run ./cmd/quill
LINT_REQUIRED_ARGS = --mode required
LINT_FULL_ARGS = --mode all --strict-recommendations --verbose

.PHONY: all build clean lint lint-required lint-fix style-install style-doctor style-coverage test

all: lint test
build:
	@mkdir -p bin
	@go build -trimpath -o bin/quill ./cmd/quill

clean:
	@GOMODCACHE="$(CURDIR)/.cache/quill/cache/go-mod" go clean -modcache
	@rm -rf -- bin .cache/quill

lint:
	@$(QUILL_CMD) check $(LINT_FULL_ARGS)

lint-required:
	@$(QUILL_CMD) check $(LINT_REQUIRED_ARGS)

lint-fix:
	@$(QUILL_CMD) fix --scope all

style-install:
	@$(QUILL_CMD) install

style-doctor:
	@$(QUILL_CMD) doctor

style-coverage:
	@$(QUILL_CMD) coverage

test:
	@go test ./...
