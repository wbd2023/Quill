# Quill product specification

## Product

Quill turns a repository's `STYLE.md` and `quill.toml` Profile into executable
style checks, safe fixes, tool installation, lock generation, and requirement
coverage.

Quill is a local developer and automation tool. It operates on a repository
available to the caller and may execute tools selected by that repository's
Profile. It is not a sandbox and is not a remote execution service.

## Release target status

[mvp.md](mvp.md) defines the ideal first public release. This
document states that release's public contract. Until the first tag is
published, release-distribution statements are commitments under validation,
not claims that archives are already available.

## Supported integration model

Quill is language-neutral through its executable, not through native libraries.
A caller pins a Quill release, invokes `quill` for a local repository, consumes
its documented output, and handles its documented exit status.

This works from Go, TypeScript, Python, Rust, Java, shell scripts, CI systems,
and any environment that can execute the supported artefact. A consuming Go
repository may use `go tool quill`; this is an executable integration, not an
importable Go API.

The supported public interfaces are:

- the `quill` command-line interface;
- versioned JSON output for machine-integrable operations;
- documented text output for people;
- documented exit status and stderr behaviour;
- `STYLE.md`, `quill.toml`, and `quill.lock`; and
- verified release binaries for supported platforms.

## Machine contract

A machine-integrable command must provide:

- an explicit `--repo-root` path for callers that cannot rely on the current
  directory;
- `--format json` output with a documented schema version;
- exactly one JSON document on stdout, with no progress or human text mixed in;
- structured operation and preparation errors in the JSON result;
- documented 0, 1, and 2 exit status semantics independent of JSON parsing;
- stderr reserved for diagnostics outside the JSON result; and
- documented signal cancellation and child-process cleanup behaviour.

The JSON protocol is an intentionally separate public DTO contract. It must not
mirror Quill's private Go structures automatically.

## Distribution

The ideal MVP release publishes checked binary archives for:

- Linux amd64 and arm64;
- macOS amd64 and arm64; and
- Windows amd64.

Each release contains one archive per supported platform and a SHA-256 checksum
file covering every archive.

Executable support and Pack or Tool capability are distinct commitments:

- The `quill` executable is supported on every listed platform. The release
  gate builds each archive on its native operating-system and architecture
  runner and executes it there, smoke-testing `quill help`, `quill version`, and
  one repository check before publishing. A platform is not supported merely
  because Quill cross-compiles for it.
- Pack and Tool capability is platform-qualified. A Pack check that requires an
  external Tool is only validated where that Tool's installer is supported. The
  cross-platform release smoke test deliberately uses a Pack check that needs no
  external Tool, so it never claims a Tool is validated on a platform where the
  Tool cannot run (for example, not every declared Tool is installable on
  Windows).

A caller outside the Go ecosystem does not need a Go toolchain to run a
published archive.

## Non-goals

Quill does not support:

- importing a Go package as a stable product interface;
- native language SDKs or FFI bindings;
- a local daemon or remote service API; or
- remote execution against repositories unavailable on the caller's machine.

## Current implementation status

The root Go facade has been removed. `internal/cli` now calls `internal/engine`
directly, and no production Go package exists at the repository root.

The source implements the CLI protocol, binary distribution workflow, and
contract tests required by ADR 0004. They remain release gates until the first
standalone release is published.
