# stylecheck

`stylecheck` is the AST-based checker used by `tools/scripts/check-style.sh` (Tier 3).

## Layout

- `main.go`
  - CLI entrypoint only.
- `analysis.go`
  - Directory walking, per-file analysis dispatch, cross-file analysis, and reporting.
- `analysis_state.go`
  - Shared analysis state and cross-file metadata types.
- `ast_type_helpers.go`
  - Shared AST and type-string helper functions.
- `scope_helpers.go`
  - Shared path-scope helpers for rule filtering.
- `text_helpers.go`
  - Shared text-style helpers used by multiple rules.
- `*_rules.go`
  - Rule-family implementations grouped by concern.
- `*_collection.go`
  - Cross-file method and implementation metadata collection helpers.
- `*_test.go`
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
- `2.2` Exported types under `internal/core/services` must end with
  `Service`, `UseCase`, or `Config` (excluding `accountref` helper package).
- `2.5` Interface methods in `internal/core/ports` must follow CRUD-L ordering.
- `2.5` Mock method order must match corresponding interface order exactly.
- `2.5` Implementation method order must match corresponding interface order exactly when
  compile-time assertions are present.
- `2.7` Parameter ordering checks: `ctx` first, secret parameters last.
- `2.8` Constructor parameter category ordering:
  repository -> service -> adapter -> config -> secret.
- `2.9` Objective file-structure ordering for top-level declarations
  (constants -> errors -> types -> assertions).

## Run

From repo root:

```bash
cd tools/stylecheck
go run . ../../internal ../../cmd ../../tests
```
