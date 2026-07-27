# ADR 0004: CLI-first language-neutral product

## Status

Accepted.

Supersedes [ADR 0003](0003-public-go-library-and-cli.md).

## Context

Quill executes repository-owned style policy against local files and invokes
repository-selected tools. It must be usable from applications and automation
written in more than one programming language.

A public Go package only serves Go callers. Keeping it as a supported product
interface would add a second compatibility contract without making Quill usable
from other runtimes. Native language bindings and FFI would multiply packaging,
error, cancellation, and compatibility work while still requiring access to the
local repository and executable tools.

A network service is not justified. It would introduce repository transfer or
mount semantics, authentication, isolation, quotas, and a new security boundary
for a product whose work belongs on the caller's machine.

## Decision

Quill is a CLI-first, language-neutral product. Its supported public interfaces
are:

- the pinned `quill` executable, including documented arguments, exit codes,
  text output, machine-readable JSON output, stderr rules, and cancellation
  behaviour;
- the versioned `STYLE.md`, `quill.toml`, and `quill.lock` repository formats;
  and
- verified release artefacts for each supported operating-system and
  architecture combination.

Applications in every language integrate by invoking the pinned executable
against a local repository. Machine-oriented commands use documented,
versioned JSON on stdout. Human-oriented text remains separate from that output.

Quill does not expose a supported in-process Go library, native SDKs, FFI
bindings, or a service API. The former temporary root `package quill` facade
has been removed. All production implementation packages remain under
`internal/`.

## Consequences

- Go, TypeScript, Python, Rust, Java, shell, and CI callers share one supported
  integration boundary.
- The command protocol becomes a compatibility contract and requires
  contract tests, schema/versioning rules, and migration notes.
- Every machine-integrable operation must define JSON success and error output,
  stdout/stderr discipline, exit semantics, and cancellation behaviour.
- Quill releases must publish an explicit platform support matrix, archives,
  checksums, and installation guidance so non-Go callers do not require a Go
  toolchain.
- `go tool quill` may remain an installation and invocation mechanism for Go
  repositories. It is not a Go library API.
- Removing the root facade deletes the duplicate public-result mapping and lets
  `internal/cli` call `internal/engine` directly.

## Alternatives considered

### Keep the public Go library and add a CLI protocol

Rejected. It creates two supported APIs with different type systems and error
models. The Go API adds no capability for non-Go consumers.

### Publish language-specific SDKs or FFI bindings

Rejected. They would duplicate compatibility and distribution work without
removing the need to access Quill's local executable-tool environment.

### Run Quill as a local or remote service

Rejected. No remote-execution requirement exists, and a service would create a
large unjustified security and operations boundary.
