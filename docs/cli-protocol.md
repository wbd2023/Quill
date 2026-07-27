# Quill CLI protocol

Version 1 is implemented in source and defended by command-protocol tests. It
becomes a supported released contract with the first tagged release that
publishes it.

The current source writes `schema_version`, `command`, and `status`. The ideal
MVP additionally requires `quill_version` on every envelope before the first
tagged release. Until then, `quill version` identifies the local executable.

## Invocation

An application invokes a pinned `quill` executable against a local repository.
It should pass `--repo-root PATH` rather than relying on current-directory
repository discovery.

Commands intended for automation accept `--format json`. Machine mode writes
one JSON document to stdout. It never writes progress, human-readable results,
or diagnostics to stdout.

## Envelope

Every JSON response uses this top-level envelope:

```json
{
  "schema_version": 1,
  "command": "check",
  "status": "ok",
  "result": {}
}
```

`schema_version` identifies the JSON protocol version. `command` is the
executed command name. `status` is `ok` or `error`.

An error response has this shape:

```json
{
  "schema_version": 1,
  "command": "check",
  "status": "error",
  "error": {
    "code": "invalid_argument",
    "message": "--scope must name a configured scope"
  }
}
```

The error `code` is a stable machine category. The message is descriptive for
people and must not be parsed for control flow.

## Command results

`check`, `doctor`, `coverage`, `fix`, `install`, and `lock` support machine
mode. Their result objects are explicit report DTOs and are defended by decoded
contract tests. Per-command public schema, field, path, and position references
remain a required ideal-MVP release deliverable.

A command can return `status: "ok"` and a nonzero exit status when its work
completed but found policy failures or an invalid toolchain. Callers must use
both the JSON result and documented exit status.

## Streams and exit status

- stdout contains only the JSON envelope in machine mode.
- stderr contains human diagnostics that are not represented by the envelope.
- exit status `0` means the command completed without blocking findings.
- exit status `1` means the command completed with blocking findings or failed
  during execution or preparation.
- exit status `2` means invalid command arguments or usage.

## Cancellation

The executable derives an operation context from `SIGINT` and `SIGTERM` and
passes it to the engine and child tools. A cancelled operation emits a JSON
error envelope in machine mode and terminates with the operating system's
signal-derived exit status. The implementation must terminate child processes,
not merely stop waiting for them.

## Compatibility

A protocol version is immutable after release. Additive result fields are
permitted within the same version. Removing or changing a field's meaning,
envelope shape, error code, stream rule, or exit semantics requires a new
protocol version and migration notes.
