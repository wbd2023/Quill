# External Pack protocol

## Status

This document defines the intended first-release external Pack interface. Until
Quill's first tagged release, its schemas and field names remain pre-release
commitments. A released protocol version is immutable except for additive
fields.

A local external Pack is a repository-contained capability. It extends Quill
through a versioned `pack.toml` manifest and a subprocess protocol, not Go
imports, native plugins, a daemon, or a remote service.

Use an External Pack only for a capability that is specific to the consuming repository or
organisation. Reusable capabilities belong in a Shipped Pack; repository-specific variation of an
existing capability belongs in that Pack's Profile-supplied Pack Policy. See
[`architecture.md`](architecture.md#pack-model-and-selection) for the full selection model.

## Source declaration and layout

A Profile declares each local source relative to the repository root:

```toml
[[pack_sources]]
path = ".quill/packs/company"
```

Quill resolves the source and its manifest through symlinks. Both must remain
inside the repository. The runtime command must resolve to an executable regular
file inside the Pack directory. Quill rejects an invalid source, manifest, or
runtime before it launches an external Pack process.

```text
.quill/packs/company/
|-- pack.toml
|-- README.md
`-- bin/
    `-- company-quill
```

## `pack.toml` schema version 1

The manifest is strict: unknown fields, unsupported versions, invalid durations,
blank or reserved Pack IDs, blank Pack names or runtime commands, duplicate Rule
IDs, and Rules with blank IDs or checks are rejected. `enabled` is reserved by
the Profile's `[packs]` table and cannot name a Pack. `runtime.timeout` defaults
to `30s` and must be positive when supplied. External fixes are not supported in
version 1, so `supports_fix` must be `false` or omitted.

```toml
schema_version = 1

[pack]
id = "company"
name = "Company Engineering Policy"
version = "0.1.0"
quill_protocol = "quill-pack-v1"

[runtime]
command = "bin/company-quill"
timeout = "30s"

[[rules]]
id = "company/no-direct-database-access"
name = "No direct database access"
group = "architecture"
check = "no-direct-database-access"
file_set = "go"
supports_fix = false
```

`pack.version`, `rules.name`, `rules.group`, and `rules.file_set` are optional.
An omitted Rule name defaults to its Rule ID; an omitted group is reported as
`external`. Rule IDs are globally unique across the complete Pack catalogue;
they are not required to use a Pack-ID prefix. `pack.quill_protocol` must be
`quill-pack-v1`.

## `quill-pack-v1` invocation

For every selected external Rule, Quill starts the declared runtime with the
repository root as its working directory. It writes exactly one JSON request to
standard input. The Pack writes JSON Lines records to standard output. Standard
error is reserved for human-readable operational messages and is not a channel
for Diagnostics.

```json
{
  "protocol": "quill-pack-v1",
  "operation": "check",
  "repository_root": "/absolute/path/to/repository",
  "pack_id": "company",
  "rule_id": "company/no-direct-database-access",
  "check_id": "no-direct-database-access",
  "scope": "all",
  "files": ["internal/service/users.go"],
  "policy": {}
}
```

`files` contains repository-relative slash paths and is always an array.
`policy` is the raw repository Profile's Pack Policy object for `pack_id`; it is
an empty object when no Pack Policy exists or when that policy is empty. Quill
passes it through without interpreting it. An external Pack therefore owns the
validation of its policy values.

The request's Pack ID, Rule ID, scope, and file selection are Quill-owned. The
runtime must not reinterpret them as authority to select another Rule, change
enforcement, or escape its declared source directory.

## Response records

A runtime writes zero or more Diagnostic records, then exactly one completion
record. Empty lines are ignored. Any other record order, unknown type, malformed
JSON, missing completion, or more than one completion is a protocol failure.

```json
{
  "type": "diagnostic",
  "code": "direct-database-access",
  "message": "Use the repository interface instead of accessing the database directly.",
  "file": "internal/service/users.go",
  "start": {"line": 42, "column": 5},
  "end": {"line": 42, "column": 18},
  "help_url": ""
}
```

A Diagnostic is one structured policy finding. `message` is required. `code` and
`help_url` are optional. `file` is either empty for a repository-level finding
or a clean repository-relative slash path. A non-empty file may have an unknown
range. A range without a file is invalid.

Lines and columns are one-based. Columns count UTF-8 bytes, not Unicode code
points or UTF-16 code units. A zero line means unknown; a zero column with a
known line means that only the column is unknown. Ranges are half-open: `start`
is inclusive and `end` is exclusive. A known end must not precede its start.

```json
{"type": "complete", "success": true}
```

A Pack can report an expected failure with a human-readable error:

```json
{"type": "complete", "success": false, "error": "policy field allowed_packages is required"}
```

Quill treats a false completion, a non-zero runtime exit, a timeout, or output
that exceeds its bound as a Rule execution error rather than a Diagnostic.

## Trust and compatibility

An external Pack executable runs with the invoking user's operating-system
permissions. Quill is not a sandbox. Review the Pack manifest, executable, and
Profile Pack Policy before running Quill on an untrusted checkout.

`pack.toml` schema versions and `quill-pack-v1` are public compatibility
contracts once released. Additive fields are allowed only when old consumers can
ignore them. Removing a field, changing its meaning, changing stream rules, or
changing Diagnostic coordinates requires a new version and migration notes.
