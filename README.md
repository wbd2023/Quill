# Tools Overview

This directory contains project tooling for build helpers and style enforcement.

## Structure

- `tools/make/`
  - Makefile include modules for grouped targets.
  - Current style target definitions live in `tools/make/style.mk`.
- `tools/scripts/`
  - Shell entrypoints, test files, and shared script infrastructure.
  - Primary entrypoints:
    - `check-style.sh`: master STYLE.md pipeline
    - `install-style-tools.sh`: installs required third-party tools
- `tools/scripts/checks/`
  - Individual shell check implementations used by `check-style.sh`.
  - Organised by concern:
    - `bash/`: Bash-specific checks
    - `go/`: Go-specific text or runner-backed checks
    - `general/`: repository-wide checks spanning multiple file types
- `tools/scripts/lib/`
  - Shared shell libraries and registry data used by style scripts.
  - `style-common.sh`: argument, scope, and common helper functions.
  - `style-registry-constants.sh`: shared registry tier and runner constants.
  - `style-registry.sh`: table-driven registry loader.
  - `style-registry.table`: single source of truth for check registration order.
- `tools/stylecheck/`
  - AST-based Go checker binary source for semantic Go rules.
  - Organised by rule family and shared analysis helpers to keep individual files focused.
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
