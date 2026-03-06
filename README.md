# Tools Overview

This directory contains project tooling for build helpers and style enforcement.

## Structure

- `tools/make/`
  - Makefile include modules for grouped targets.
  - Current style target definitions live in `tools/make/style.mk`.
- `tools/style/`
  - Style tooling domain root.
  - `entrypoints/`: runnable style commands (`check-style.sh`, `install-tools.sh`)
  - `checks/`: shell checks grouped by concern (`bash`, `go`, `repository`)
  - `internal/`: shared shell helpers and registry definitions
  - `tests/`: Go black-box tests for style tooling
  - `ast/`: AST-based Go checker module
  - See `tools/style/README.md` for style tooling details.
- `tools/style/checks/`
  - Individual shell check implementations used by `check-style.sh`.
  - Organised by concern:
    - `bash/`: Bash-specific checks
    - `go/`: Go-specific text or runner-backed checks
    - `repository/`: repository-wide checks spanning multiple file types
- `tools/style/internal/`
  - Shared shell libraries and registry data used by style scripts.
  - `common.sh`: argument, scope, and common helper functions.
  - `registry-constants.sh`: shared registry tier and runner constants.
  - `registry.sh`: table-driven registry loader.
  - `runner.sh`: shared execution and reporting helpers for `check-style.sh`.
  - `registry.table`: single source of truth for check registration order.

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

- Shell tooling tests: `go test ./tools/style/...`
- AST checker tests: `(cd tools/style/ast && go test ./...)`
