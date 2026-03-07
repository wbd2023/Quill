# Stylecheck

`stylecheck` is the Go AST-based checker used by
`tools/style/entrypoints/check-style.sh` (Tier 3).

## Layout

- `main.go`
  - CLI entrypoint only.
- `internal/lint/check.go`
  - CLI argument parsing and top-level lint orchestration.
- `internal/lint/scan.go`
  - Directory walk, file processing, cross-file indexing, and reporting.
- `internal/lint/report/violation.go`
  - Shared violation type used across checker packages.
- `internal/lint/paths/paths.go`
  - Path-scope constants and predicates.
- `internal/lint/syntax/syntax.go`
  - Syntax-tree helpers shared by rules and index passes.
- `internal/lint/rules/`
  - Rule-family implementations plus small text helpers used by rule checks.
- `internal/lint/index/`
  - `types.go` holds the shared declaration models.
  - `interfaces.go` and `implementations.go` build cross-file indexes.
  - `order.go` validates ordering for interfaces, mocks, and implementations.
- `internal/lint/tests/`
  - Rule-family regression tests plus shared test helpers.

## Implemented checks

- `2.2` Named return values are required.
- `2.2` Naked returns are disallowed when return values are named.
- `2.2` Parameter type elision is disallowed.
- `2.1` Error-message context for `fmt.Errorf` and `errors.New` must be lower-case and
  must not end with punctuation.
- `2.1` `fmt.Errorf` arguments must not include secret-like identifiers.
- `2.1` Sentinel errors (`Err...`) are restricted to `internal/core/domain/errors.go`.
- `2.1` Adapters must not propagate bare `return err` / `return ..., err` forms.
- `2.3` Trailing inline comments must start lower-case and avoid ending punctuation.
- `2.2` Direct casts to key domain ID aliases are disallowed in all Go code
  (including tests) outside `internal/core/domain`; parser/constructor helpers are required.
  This check is type-aware and also catches alias-import and local type-alias bypasses.
- `2.2` Single-letter variable names are restricted (`i`, `j`, `k`, `_` only).
- `2.2` Exported types under `internal/*/application/service` must end with
  `Service` or `Config`.
- `2.5` Interface methods in `internal/*/application/port` must follow CRUD-L ordering.
- `2.5` Mock method order must match corresponding interface order exactly.
- `2.5` Implementation method order must match corresponding interface order exactly when
  compile-time assertions are present.
- `2.7` Parameter ordering checks: `ctx` first, secret parameters last.
- `2.8` Constructor parameter category ordering:
  repository -> service -> adapter -> config -> secret.
- `2.9` Objective file-structure ordering for top-level declarations:
  constants -> errors -> types -> assertions.

## Run

From repo root:

```bash
cd tools/style/stylecheck
go run . ../../internal ../../cmd ../../test
```
