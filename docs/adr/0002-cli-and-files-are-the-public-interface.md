# ADR 0002: CLI and files are the public interface

## Status

Accepted.

## Context

Quill must be reusable across repositories, but its current consumers need command execution and
repository-owned policy files rather than an in-process Go library. Publishing packages early would
freeze vocabulary, package boundaries, and orchestration APIs before a concrete caller establishes
which parts require compatibility.

The implementation is a modular monolith. Internal packages separate Profile compilation, Pack
composition, execution, tool installation, Checks, and reporting, while one engine facade
coordinates a complete operation.

## Decision

Quill exposes these supported integration surfaces:

- the `quill` command and its documented exit codes and text/JSON formats;
- `STYLE.md` requirement metadata;
- the `quill.toml` Profile format;
- the `quill.lock` archive-hash format.

All Go implementation packages remain under `internal/`. Quill will not add `pkg/`, public aliases,
re-exports, or compatibility wrappers without a concrete external Go consumer and an explicit API
stability decision.

Consuming Go 1.24 repositories should pin the command with a `tool` directive and invoke it
through `go tool quill`. Other distribution mechanisms must provide equivalent version pinning and
integrity.

## Consequences

- Package boundaries may be improved without breaking downstream Go builds.
- CLI and file-format changes require migration analysis, tests, and release notes.
- Embedding Quill in another Go process is intentionally unsupported today.
- A future public API requires a new ADR identifying its consumer, compatibility promise, and the
  smallest coherent package surface.

## Alternatives considered

### Publish selected internal packages now

Rejected. Candidate packages expose coupled domain types and no current consumer needs them. Moving
them would create a compatibility burden without proving the API shape.

### Keep a duplicate embedded copy in each consumer

Rejected. Copies drift, hide ownership, and force consumers to carry Quill's implementation and test
suite. A pinned command preserves one release and one source of truth.

### Use a checksummed binary as the only integration

Deferred. Binary distribution needs a supported platform matrix, archive contract, and checksum
publication process. The Go 1.24 tool mechanism already gives Go consumers a pinned module version.
