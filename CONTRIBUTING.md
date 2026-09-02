# Contributing to Quill

Quill is in the design phase of a clean-slate Rust rewrite. The branch intentionally has no active
implementation or build commands yet. Discuss substantial product, compatibility, domain, or
architecture changes before adding code.

## Contributor licence agreement

Every external contribution requires acceptance of the
[Quill Individual Contributor Licence Agreement](CLA.md) before it can be merged.

Read [CONTRIBUTOR_PRIVACY.md](CONTRIBUTOR_PRIVACY.md) before providing the personal information
used to record acceptance. A contributor may withdraw consent before acceptance is recorded by
contacting the maintainer. Recorded acceptance and contribution records are retained as described
in that notice.

If you are under 18, or if an employer, university, client, funding body, or other person may own
your contribution, do not use the standard individual acceptance flow until authority to contribute
has been confirmed.

## Contribution provenance

Disclose in the pull request:

- third-party code, text, data, fixtures, generated material, or ideas incorporated into the change;
- the source, copyright owner, and licence of that material;
- any use of code-generation or AI systems that may affect authorship or licensing; and
- any employer, university, client, funding, confidentiality, or export-control obligation that may
  affect your authority to contribute.

Do not submit confidential material, personal data without authority, secrets, credentials, or
content whose licence is incompatible with Apache-2.0 distribution.

## Current change standard

Until the Rust product architecture and build surface are approved:

- prefer product-contract, feature-inventory, domain-model, and architecture work;
- treat retained Go-era documents as evidence rather than automatic requirements;
- do not add placeholder crates, modules, abstractions, compatibility layers, or generated output;
- keep changes independently reviewable and explain the concrete Quill behavior behind each major
  abstraction; and
- state exactly what was verified, because no repository-wide Rust command exists yet.

Security-sensitive changes must preserve or deliberately replace the relevant requirements in
[SECURITY.md](SECURITY.md).
