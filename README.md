# Quill

Quill is a reproducible, extensible style-policy platform for software repositories. A repository
owns written engineering requirements and machine-readable policy; Quill validates that policy,
runs its selected checks and fixes, and reports how automation covers the written requirements.

## Rewrite status

This branch begins a clean-slate Rust redesign. It intentionally contains no active implementation,
build definition, command-line program, self-hosting Profile, lockfile, CI workflow, or release
workflow yet.

The final Go implementation is preserved by this branch's parent commit:

```text
3ed482e569b92cd6b4b7f1be5a0b80d64fbaa4e5
```

The source-only `v0.2.0` release preserves the final Go implementation and its product/protocol
documents as the rewrite's archaeology anchor. Publication is now established; adoption and the
Rust compatibility commitment still require evidence rather than assumption.

## Design intent

The Rust implementation will be designed from Quill's product requirements and difficult edge
cases, not by translating Go packages mechanically. Before adding a Cargo workspace or source
modules, the rewrite will classify current capabilities as retained, redesigned, or removed and
will establish the domain and execution model.

Retained documents are design and compatibility evidence. They are not automatically authoritative
specifications for the Rust implementation.

## Reference material

- [STYLE.md](STYLE.md) is the existing human policy and an important parser/policy corpus. Its
  Go-specific requirements require deliberate Rust revision.
- [docs/product.md](docs/product.md) records the current product direction.
- [docs/cli-protocol.md](docs/cli-protocol.md) records the final-Go machine interface.
- [docs/pack-protocol.md](docs/pack-protocol.md) records the final-Go External Pack interface.
- [ADR 0004](docs/adr/0004-cli-first-language-neutral-product.md)
  records the CLI-first product decision.
- [ADR 0006](docs/adr/0006-external-pack-protocol-is-a-public-interface.md)
  records the subprocess extension decision.

The removed Go implementation, tests, build machinery, self-hosting configuration, architecture
documents, and historical ADRs remain available through Git history.

## Governance

See [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md),
[CLA.md](CLA.md), and [CONTRIBUTOR_PRIVACY.md](CONTRIBUTOR_PRIVACY.md).

Quill is licensed under the [Apache License, Version 2.0](LICENSE). See [NOTICE](NOTICE) for
attribution.
