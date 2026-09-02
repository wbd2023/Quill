# Quill product specification

> **Rust rewrite status:** This document is retained as product and compatibility evidence from the
> final Go implementation. It is not automatically authoritative for the Rust redesign.

## Product

Quill turns a repository's `STYLE.md` and `quill.toml` Profile into executable
style checks, safe fixes, tool installation, lock generation, and requirement
coverage.

Quill is a local developer and automation tool. It operates on a repository
available to the caller and may execute tools selected by that repository's
Profile. It is not a sandbox and is not a remote execution service.

## Release target status

The repository has a `v0.1.0` tag. This document also describes final-Go capabilities added after
that tag, including interfaces whose publication and adoption have not been established. The Rust
rewrite must classify each surface as retained, redesigned, or removed before assigning a new
compatibility promise.

## Supported integration model

Quill is language-neutral through its executable, not through native libraries.
A caller pins a Quill release, invokes `quill` for a local repository, consumes
its documented output, and handles its documented exit status.

This model works from Go, TypeScript, Python, Rust, Java, shell scripts, CI systems, and any
environment that can execute a supported artefact. The Rust rewrite must preserve or deliberately
replace executable integration without accidentally creating a language-specific library contract.

The supported public interfaces are:

- the `quill` command-line interface;
- versioned JSON output for machine-integrable operations;
- documented text output for people;
- documented exit status and stderr behaviour;
- `STYLE.md`, `quill.toml`, and `quill.lock`;
- local external-Pack `pack.toml` manifests and the versioned Pack subprocess
  protocol; and
- verified release binaries for supported platforms.

## Repository formats and Pack extensions

`quill.toml`, `quill.lock`, and local external-Pack `pack.toml` have explicit
schema versions. `STYLE.md` requirement metadata has no independent schema
marker: its accepted grammar and compatibility are part of the pinned Quill
release contract. Additive metadata must state the minimum Quill version; an
incompatible metadata change requires a documented migration.

[pack-protocol.md](pack-protocol.md) defines the external Pack manifest,
subprocess records, Diagnostic coordinates, compatibility rules, and trust
model. The external Pack interface is local and repository-contained; it is
not a remote registry, Go SDK, native plugin, or service API.

## Machine contract

A machine-integrable command must provide:

- an explicit `--repository-root` path for callers that cannot rely on the current
  directory;
- `--format json` output with a documented schema version;
- exactly one JSON document on stdout, with no progress or human text mixed in;
- structured operation and preparation errors in the JSON result;
- documented 0, 1, and 2 exit status semantics independent of JSON parsing;
- stderr reserved for human-readable operational messages outside the JSON
  result; and
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

- importing a language-specific library as a stable product interface;
- native language SDKs or FFI bindings;
- a local daemon or remote service API; or
- remote execution against repositories unavailable on the caller's machine.

## Current implementation status

This branch contains no active implementation. The final Go implementation is preserved at parent
commit `3ed482e569b92cd6b4b7f1be5a0b80d64fbaa4e5`. Its CLI, distribution workflow, local External
Pack interface, and contract tests remain evidence for the Rust product decision inventory.
