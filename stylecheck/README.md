# stylecheck

`stylecheck` is the AST-based checker used by `tools/scripts/check-style.sh` (Tier 3).

## Implemented checks

- `2.2` Named return values are required.
- `2.2` Naked returns are disallowed when return values are named.
- `2.2` Parameter type elision is disallowed.
- `2.2` Direct casts to key domain ID aliases are disallowed outside
  `internal/core/domain`; parser/constructor helpers are required.
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
