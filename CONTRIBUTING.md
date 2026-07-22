# Contributing to Quill

Quill is a modular Go CLI. Contributions should improve a complete user-visible contract without
introducing speculative extension points or public packages.

Bug reports, feature requests, design discussion, and documentation feedback are welcome. Please
discuss substantial changes with the maintainer before investing significant implementation work.

## Development setup

Required tools:

- Go 1.24.5
- Node.js 20
- GNU Make

Install the repository-pinned tools:

```sh
make style-install
```

The installer writes tools and caches beneath the ignored `.cache/quill/` directory.

## Contributor licence agreement

Every external contribution requires acceptance of the
[Quill Individual Contributor Licence Agreement](CLA.md) before it can be merged.

Quill uses [CLA Assistant][cla-assistant] to authenticate contributors, record acceptance, and
report CLA status on pull requests. The automated CLA status check must pass before merging.

For your first contribution:

1. Open the pull request using the GitHub account that submitted the Contribution.
2. Follow the signing link posted by CLA Assistant.
3. Review the signing copy of the CLA and its exact Gist revision.
4. Complete the private contributor fields and declarations.
5. Accept the Agreement through CLA Assistant.

Your full legal name and declaration responses are not posted to the Quill repository or pull
request by Quill. CLA Assistant and GitHub process the information described in
[CONTRIBUTOR_PRIVACY.md](CONTRIBUTOR_PRIVACY.md).

The CLA is individual-only. If an employer, company, university, client, or other legal entity owns
or may own the Contribution, contact the maintainer before submitting it. Written authorisation or
a separate entity contributor agreement may be required.

If you are under 18, contact the maintainer before contributing. Do not use the standard
individual acceptance process.

A materially revised CLA requires fresh acceptance. A pull-request checkbox or Developer
Certificate of Origin sign-off does not replace the CLA process.

General discussion, ideas, feature requests, and bug reports are not Contributions unless specific
material is expressly submitted for inclusion in Quill.

## Contribution provenance

Disclose in the pull request:

- all third-party material, including its source, copyright owner, and licence;
- any material whose authorship, ownership, or licensing status may be uncertain, including its
  source and applicable terms where known; and
- any employer, university, client, funding, confidentiality, or other obligation that may affect
  your authority to contribute.

You remain responsible for reviewing the complete Contribution and confirming that you can grant
the rights stated in the CLA. Do not submit confidential information, trade secrets, personal
information, or material copied from another project unless its licence and inclusion have been
expressly disclosed and are compatible with Quill.

## Before submitting a change

Run:

```sh
make lint
make test
make build
```

For focused work, run the smallest affected package test first, then the complete commands above.
Tests must be deterministic and must not depend on a parent checkout, a global tool installation,
or network access unless they are explicitly isolated integration tests.

## Design rules

- Keep one CLI entrypoint in `cmd/quill`.
- Keep implementation packages under `internal/`; do not add `pkg/` without a real external
  consumer.
- Keep Profile policy, TOML persistence, Pack declarations, concrete Checks, execution, and
  reporting as separate responsibilities.
- Keep project-specific policy in `quill.toml`, not in shipped Pack code.
- Do not use global registration or `init` side effects.
- Treat downloads, archives, checksums, process execution, and filesystem writes as security
  boundaries.
- Preserve text and JSON output contracts and exit-code semantics.

Architecture tests in `internal/architecture` enforce important import boundaries. Update an
architecture test only when the intended ownership boundary has genuinely changed.

## Style and documentation

`STYLE.md` is normative and `quill.toml` selects its executable Rules. Use Australian English,
ASCII punctuation, and the configured 100-column line limit. Exported Go declarations require
useful doc comments. Comments should explain contracts or non-obvious decisions rather than repeat
syntax.

Update user documentation when a CLI flag, command, output field, Profile key, lockfile field, or
security boundary changes.

## Tests

Tests should defend observable behaviour, boundaries, invariants, transitions, precedence, or real
error paths. Repository-shaped fixtures belong under temporary directories or package `testdata`.
Do not couple tests to an enclosing repository.

Security-sensitive installer changes require regression coverage for the affected rejection path,
such as checksum mismatch, path traversal, links, archive size limits, or unsafe replacement.

## Commit and pull request scope

Keep changes coherent and reviewable. A package move must update every import, architecture rule,
test, command, and document in the same change. Do not leave compatibility aliases in this pre-1.0
repository unless an external consumer requires one.

Pull requests should state:

- the user-visible problem and final behaviour;
- affected CLI, Profile, or security contracts;
- verification commands and observed results;
- any migration required for consuming repositories;
- all third-party material and any material with uncertain authorship, ownership, or licensing
  status; and
- whether a separate contribution arrangement may be required.

[cla-assistant]: https://cla-assistant.io/
