# Quill documentation

This directory holds Quill's durable product, protocol, architecture, planning,
and decision records. Each document has one owner and purpose.

## Product and delivery

- [product.md](product.md) defines Quill's supported public integration model and
  release contract.
- [mvp.md](mvp.md) defines the ideal first-release outcome, non-goals, and
  acceptance criteria.
- [roadmap.md](roadmap.md) is the mutable, ordered delivery plan for that MVP.

## Interfaces and architecture

- [cli-protocol.md](cli-protocol.md) defines the public machine CLI contract.
- [pack-protocol.md](pack-protocol.md) defines the public local external-Pack
  manifest and subprocess protocol.
- [architecture.md](architecture.md) describes current package ownership,
  dependency direction, and runtime flow.
- [adr/](adr/) records durable architectural decisions and their alternatives.

## Engineering review

- [codebase-review.md](codebase-review.md) is the mutable, design-first
  whole-codebase review programme. It owns review scope, sequence, evidence,
  and completion state without redefining product, protocol, architecture,
  security, or delivery contracts.

## Repository-root documents

- [../README.md](../README.md) is the project landing page, installation guide,
  and human CLI overview.
- [../CONTRIBUTING.md](../CONTRIBUTING.md) explains how to contribute.
- [../SECURITY.md](../SECURITY.md) defines the trust model and vulnerability
  reporting process.
- [../TODO.md](../TODO.md) is the tactical maintenance queue for small,
  repository-local work outside the ideal MVP roadmap.
- [../STYLE.md](../STYLE.md), [../quill.toml](../quill.toml), and
  [../quill.lock](../quill.lock) are the repository's executable policy inputs,
  not general documentation.
- [../CLA.md](../CLA.md), [../CONTRIBUTOR_PRIVACY.md](../CONTRIBUTOR_PRIVACY.md),
  [../LICENSE](../LICENSE), and [../NOTICE](../NOTICE) remain at the repository
  root because they are legal documents with stable public paths.
