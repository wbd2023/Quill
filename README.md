# Tools Overview

This directory contains project tooling for build helpers and style enforcement.

## Structure

- `tools/make/`
  - Makefile include modules for grouped targets.
  - Current style target definitions live in `tools/make/style.mk`.
- `tools/scripts/`
  - Shell-based checkers and installers.
  - Primary entrypoints:
    - `check-style.sh`: master STYLE.md pipeline
    - `install-style-tools.sh`: installs required third-party tools
- `tools/scripts/lib/`
  - Shared shell libraries used by style scripts.
  - `style-common.sh`: argument, scope, and common helper functions.
  - `style-registry.sh`: table-driven registry loader.
  - `style-registry.table`: single source of truth for check registration order.
- `tools/stylecheck/`
  - AST-based Go checker binary source.
  - See `tools/stylecheck/README.md` for checker-specific details.

## Recommended Usage

Run style checks through `make` targets rather than calling each script directly.

- `make style`
  - Required STYLE.md checks only.
- `make style-all`
  - Required checks plus recommendation checks.
- `make style-all-strict`
  - Same as `style-all`, but recommendation findings fail the run.

Verbose output can be enabled with `VERBOSE=true`.

- `make style VERBOSE=true`
- `make style-all VERBOSE=true`
- `make style-all-strict VERBOSE=true`

## Testing

- Shell tooling tests: `go test ./tools/scripts/...`
- AST checker tests: `(cd tools/stylecheck && go test ./...)`
