# Quill CLI protocol

This document records version 1 as implemented on final-Go `main`. It was added after the
repository's `v0.1.0` tag, so its publication and adoption status must be established before the
Rust rewrite preserves it as version 1 or replaces it with an incompatible contract.

## Invocation

An application invokes a pinned `quill` executable against a local repository.
It should pass `--repository-root PATH` rather than relying on current-directory
repository discovery.

Commands intended for automation accept `--format json`. Machine mode writes
one JSON document to stdout. It never writes progress, human-readable results,
or human-readable operational messages to stdout.

## Envelope

Every JSON response uses this top-level envelope:

```json
{
  "schema_version": 1,
  "quill_version": "v1.2.3",
  "command": "check",
  "status": "ok",
  "result": {}
}
```

`schema_version` identifies the JSON protocol version. `quill_version` is the
version of the executable that produced the response. `command` is the executed
command name. `status` is `ok` or `error`.

An error response has this shape:

```json
{
  "schema_version": 1,
  "quill_version": "v1.2.3",
  "command": "check",
  "status": "error",
  "error": {
    "code": "invalid_argument",
    "message": "--scope must name a configured scope"
  }
}
```

An `ok` envelope contains `result` and no `error`. An `error` envelope contains
`error` and no `result`. The error `code` is a stable machine category. The
message is descriptive for people and must not be parsed for control flow.

## Command results

`check`, `doctor`, `coverage`, `fix`, `install`, `lock`, `list`, and `explain`
support machine mode. The field paths below are relative to the envelope's
`result`. All arrays use `[]`, and `?` marks an omitted field.

### `check`

- `result.result.entries[]`: `rule_id`, `name`, `group`, `enforcement`, `scope`,
  `status`, `requirements[]`, `diagnostics[]`, `execution_error?`, `command?`.
- `result.result.entries[].diagnostics[]`: `code`, `file?`, `range?`, `message`,
  `help_url?`. A Diagnostic is one structured policy finding, not an operational
  error or stderr message.
- `result.result.entries[].diagnostics[].file`: a clean repository-relative slash path,
  or omitted for a repository-level finding.
- `result.result.entries[].diagnostics[].range`: `start.line`, `start.column?`,
  `end?.line`, `end?.column`. Lines and columns are one-based; columns count
  UTF-8 bytes. A zero line is unknown, and a zero column with a known line is
  unknown. The range is half-open: start is inclusive and end is exclusive.
  A range without a file, or an end before its start, is invalid.
- `result.result.entries[].command`: `exit_code`, `timed_out`, `truncated`.
- `result.summary`: `Passed`, `Warned`, `Failed`, `Blocked`, `Skipped`, `Errored`.
- `result.groups[]`: `group` and the matching `entries[]` view.
A rule whose check ran but could not complete records `execution_error?` and a
per-entry `status` of `error`. The command still returns `status: "ok"` and
exit `1`: the work completed, and the failing rule is one of its findings.
Only a failure that prevents the command itself from completing - preparation,
execution orchestration, cancellation, or rendering - returns an `error`
envelope with code `operation_failed` (or `cancelled` for cancellation).

### `coverage`

- `report.requirements[]`: `id`, `section`, `text`, `mode`, `reason?`,
  `rule_ids[]`.
- `report.sections[]`: `section`, `title`, `status`, `requirement_count`,
  `automated_count`, `review_only_count`, `manual_deferred_count`.
- `requirement_totals`: `automated`, `review_only`, `manual_deferred`.
- `section_totals`: `automated`, `partial`, `review_only`, `manual`.
- `outstanding[]`: the same schema as `report.requirements[]`.
- `outstanding_by_mode`: a map from coverage mode to outstanding count.

### `fix`

- `scope`: selected scope.
- `toolchain.all_valid`: whether every required tool passed inspection.
- `toolchain.statuses[]`: the tool status schema below.
- `rules[]`: `rule_id`, `name`, `group`, `enforcement`, `scope`, `exit_code`,
  `timed_out`, `truncated`, `execution_error?`.

### `doctor` and `install`

Both commands return:

- `result.result.statuses[]`: `id`, `name`, `path`, `version`, `pinned_version`,
  `valid`, `issue?`.
- `result.all_valid`: whether every inspected tool is valid.

### `lock`

- `path`: absolute path of the rewritten `quill.lock`.
- `archive_count`: number of archive-tool entries written.

### `list`

Exactly one key is present:

- `packs[]`: `id`, `name`, `active`, `provenance`, `rules`, `tools`.
- `rules[]`: `id`, `pack`, `provenance`, `name`, `active`, `enforcement?`,
  `scope?`, `fix`.
- `tools[]`: `id`, `name`, `command?`, `pin?`, `packs[]`.
- `scopes[]`: `name`, `roots[]`, `default`.

The key matches the selector supplied to `quill list`.

### `explain`

- `rule`: `id`, `name`, `group`, `pack`, `binding`, `check`, `fix?`.
- `rule.pack`: `id`, `name`, `provenance`, `policy?`.
- `rule.binding`: `enforcement`, `scope`, `requirements[]`.
- `rule.check` and `rule.fix`: `category`, `tools?`, `file_set?`, `language?`.

`init` and `version` are human-facing commands and do not implement JSON mode.

## Streams and exit status

- stdout contains only the JSON envelope in machine mode.
- stderr contains human-readable operational messages that are not represented
  by the envelope.
- exit status `0` means the command completed without blocking findings.
- exit status `1` means the command completed with blocking findings or failed
  during execution or preparation.
- exit status `2` means invalid command arguments or usage.

A command can return `status: "ok"` and exit `1` when its work completed but
found policy failures, per-rule execution errors, or an invalid toolchain.
Callers must use both the JSON result and documented exit status.

## Invalid arguments and cancellation

Quill validates command grammar before repository discovery or command work.
Malformed invocations normally write a human error and contextual usage to stderr
with exit status `2`. If a command has established `--format json` before
parsing fails, Quill instead writes one `invalid_argument` envelope to stdout
and exits `2`.

Repository discovery, filesystem access, profile loading, preparation,
execution, and rendering failures are `operation_failed` with exit status `1`.
Repository-dependent selections that are syntactically valid but unavailable in
the loaded profile, such as an unknown scope or inactive rule, are
`invalid_argument` with exit status `2`. Callers must pass canonical long flags
such as `--format`; undocumented single-dash long forms are invalid.

The executable derives an operation context from `SIGINT` and `SIGTERM` and
passes it to the engine and child tools. An operation is `cancelled` only when
its returned error wraps `context.Canceled`; a timeout remains
`operation_failed`. After the first signal Quill restores the default signal
handling, so a second signal can force termination. A cancelled operation emits
a JSON error envelope in machine mode before terminating with the operating
system's signal-derived exit status. The implementation terminates child
processes, not merely stops waiting for them.

## Compatibility

A protocol version is immutable after release. Additive result fields are
permitted within the same version. Removing or changing a field's meaning,
envelope shape, error code, stream rule, or exit semantics requires a new
protocol version and migration notes.
