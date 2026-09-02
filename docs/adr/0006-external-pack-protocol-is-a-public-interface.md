# ADR 0006: External Pack protocol is a public interface

## Status

Accepted for revalidation during the Rust rewrite.

## Context

ADR 0004 establishes Quill as a CLI-first, language-neutral product. Its public
interfaces include the executable, repository policy formats, documented CLI
output, and release artefacts; Quill does not provide an in-process Go API,
SDK, FFI, daemon, or remote service.

The first-release product also requires a repository to add a local external
Pack without forking, importing, or rebuilding Quill. The implementation already
loads repository-contained Pack sources, validates strict `pack.toml` manifests,
and executes their declared runtimes through `quill-pack-v1`. Leaving those
formats undocumented would make an implemented extension mechanism an unstable
private accident and leave Pack authors without a compatibility or trust model.

## Decision

A local external Pack is a supported public extension interface. Its public
contracts are:

- the Profile `[[pack_sources]]` declaration;
- the versioned `pack.toml` manifest; and
- the versioned `quill-pack-v1` subprocess protocol.

[pack-protocol.md](../pack-protocol.md) owns their exact schemas, stream rules,
Diagnostic coordinates, compatibility policy, and Pack-author trust model.

External Packs remain repository-local in the first release. Sources and runtime
executables must resolve inside the consuming repository. Quill validates the
source, manifest, executable containment, and protocol output before trusting
those inputs, but does not sandbox a Pack executable.

The external request's `policy` field is the raw Pack Policy object from the
repository Profile. Quill forwards it without interpreting Pack-specific values;
the external Pack validates its own policy.

This decision augments ADR 0004. It neither restores a public Go API nor creates
a second orchestration path. External Pack values are data and subprocess
records, never reflections of Quill's private Go types.

## Consequences

- `pack.toml` schema versions and released Pack protocol versions are immutable
  except for compatible additive fields.
- Removing a field, changing its meaning, stream behaviour, Diagnostic
  coordinates, or execution semantics requires a new version and migration
  notes.
- Product, security, and release documentation must name external Pack manifests
  and executables as reviewable public/trust surfaces.
- Pack protocol tests are compatibility tests, not only implementation tests.
- Remote registries, dependency resolution, remote executable installation,
  native plugins, and arbitrary sources outside the repository remain out of
  scope.

## Alternatives considered

### Keep external Packs as a private implementation detail

Rejected. The first-release MVP requires externally authored local Packs, and
the source already accepts their persisted data and executes their binaries.
Treating that behaviour as private would make a compatibility promise without an
owner.

### Remove external Packs from the first release

Rejected. External adoption is the MVP proof that Quill is reusable beyond its
own repository. No product requirement supports shrinking that scope.

### Publish a Go SDK or native plugin interface

Rejected. It would violate ADR 0004 by adding a language-specific public API and
would couple Pack authors to private implementation types. A small subprocess
protocol works from any implementation language.

### Support remote Pack registries now

Rejected. Registry trust, dependency resolution, authentication, distribution,
and remote executable installation are separate product and security problems.
A repository-contained local source is the smallest useful extension interface.
