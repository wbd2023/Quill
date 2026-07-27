# ADR 0003: Public Go library and CLI

## Status

Superseded by [ADR 0004](0004-cli-first-language-neutral-product.md).

This records the historical public-library decision and its implementation.

## Context

Quill is used through its command and repository-owned policy files. Some consumers also need to
load and execute the same policy in-process, inspect structured results, and decide their own
presentation without spawning a command.

Publishing implementation packages would freeze Profile compilation, Pack composition, Driver
selection, Checks, and platform boundaries before those package shapes are stable. Maintaining a
separate library orchestration path would duplicate the CLI's operation behaviour and permit drift.

## Decision

Quill exposes one importable Go package at the module root:

```go
import "github.com/wbd2023/quill"
```

The public package exposes `Runner`, its operation methods, typed operation options, and structured
presentation-free result values. It does not expose internal Profile, Pack, Driver, Check,
toolchain, installer, or report types.

`internal/engine` remains the private application facade. The root package maps public values to and
from private implementation values. The dependency direction is:

```text
external Go caller -> quill.Runner -> internal/engine -> internal capabilities
cmd/quill -> internal/cli -> quill.Runner -> internal/engine -> internal capabilities
```

The CLI preserves its documented command names, flags, exit codes, and text and JSON output. In
particular, `quill doctor` invokes the library's `Runner.Inspect` operation.

The module path is `github.com/wbd2023/quill`. The root package is the only public production Go
package. Quill does not add `pkg/`, public aliases, re-exports of internal packages, or
compatibility wrappers.

## Consequences

- The CLI and Go callers execute the same operations and receive the same underlying outcomes.
- Public API changes require semantic-versioning analysis, documentation, examples, and migration
  notes.
- Internal package boundaries may continue to evolve without downstream Go migration as long as the
  root package contract remains stable.
- `internal/report` renders public result values for the CLI and remains private presentation code.
- The capability-oriented internal layout remains intact. Generic architecture-layer directories are
  not introduced.
- Before v1, breaking public API changes require a minor-version increment and migration notes.
  After v1, breaking changes require a new major-version module path.

## Alternatives considered

### Keep the CLI and files as the only public interface

Rejected. It prevents valid in-process integrations from obtaining structured results without
process management and output parsing.

### Publish selected internal packages

Rejected. This would expose coupled implementation vocabulary and make internal refactoring a
compatibility concern.

### Add a named public subpackage

Rejected. Quill's primary library is the product itself, so the module root gives the conventional
and shortest import path. `cmd/quill` remains reserved for the binary.

### Give the CLI its own Engine orchestration path

Rejected. Two operation paths would drift in validation, result mapping, and error behaviour.
The CLI must call the root public facade.
