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

The CLA is individual-only. If an employer, company, university, client, or other legal entity owns
or may own the contribution, contact the maintainer before submitting it. Written authorisation or
a separate entity contributor agreement may be required.

If you are under 18, contact the maintainer before contributing. Do not use the standard acceptance
statement.

For your first contribution under CLA Version 1.1, post the following statement as a comment on the
pull request from the same GitHub account that submitted it:

```text
I, FULL LEGAL NAME, have read and accept the Quill Individual Contributor Licence Agreement,
Version 1.1, at commit CLA_COMMIT_SHA, including the moral-rights consent in Section 7.

I confirm that I am at least 18 years old, that I have authority to grant its rights, and that this
acceptance applies to pull request #PULL_REQUEST_NUMBER and my later Accepted Contributions under
Version 1.1.

Signed electronically by @GITHUB_USERNAME on YYYY-MM-DD.
```

Replace every capitalised placeholder. `CLA_COMMIT_SHA` must be the full Git commit SHA of a
repository revision containing the Version 1.1 `CLA.md`. The acceptance comment is public.

A maintainer must verify the statement, apply the `cla:accepted` pull-request label, and reply
using this form before merging:

```text
CLA acceptance recorded for FULL LEGAL NAME (@GITHUB_USERNAME), Version 1.1 at
CLA_COMMIT_SHA, on YYYY-MM-DD.
```

The acceptance comment, maintainer reply, pull request, and referenced repository commit form the
acceptance record.

A Version 1.1 acceptance remains valid for later contributions from the same GitHub account. A
revised CLA requires fresh acceptance. A pull-request checkbox or Developer Certificate of Origin
sign-off does not replace the acceptance statement.

General discussion, ideas, feature requests, and bug reports are not Contributions unless specific
material is expressly submitted for inclusion in Quill.

## Contribution provenance

Disclose in the pull request:

- all third-party material, including its source, copyright owner, and licence;
- any other material whose authorship, ownership, or licensing status may be uncertain, including
  its source and applicable terms where known; and
- any employer, university, client, funding, confidentiality, or other obligation that may affect
  your authority to contribute.

You remain responsible for reviewing the complete contribution and confirming that you can grant
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
Tests must be deterministic and must not depend on a parent checkout, a global tool installation, or
network access unless they are explicitly isolated integration tests.

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
- all third-party material and any other material with uncertain authorship, ownership, or
  licensing status; and
- the CLA acceptance comment or a link to the contributor's earlier Version 1.1 acceptance.
