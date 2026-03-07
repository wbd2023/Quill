# Style Tooling

This directory contains the complete STYLE.md automation pipeline.

## Layout

- `entrypoints/`
  - `check-style.sh`: runs the multi-tier style pipeline.
  - `install-tools.sh`: installs third-party style tooling dependencies.
- `checks/`
  - `bash/`: Bash-focused style and lint checks.
  - `go/`: Go-focused text and linter-backed checks.
  - `repository/`: repository-wide checks (spelling, ASCII, markdown, naming, headers).
- `internal/`
  - Shared shell helpers (`common.sh`, `runner.sh`) and registry files
    (`registry-constants.sh`, `registry.sh`, `registry.table`).
- `tests/`
  - Go black-box tests for entrypoints, registry loading, and check behaviour.
- `stylecheck/`
  - AST-based Go checker module used by Tier 3 checks.

## Usage

- `make style`
- `make style-all`
- `make style-all-strict`

Set `VERBOSE=true` to print failing command output.
