# ADR 0005: Kong owns CLI grammar

## Status

Accepted.

## Context

Quill has ten commands with shared flags, required positionals, static enums,
and a machine-readable protocol. Its former CLI layer spread one command across
an option struct, `flag.FlagSet` construction, manual parsing and validation,
usage rendering, a mutable command registry, and a second raw-argument parse to
recover machine mode after errors.

That duplicated the command grammar, made command additions non-local, and
risked divergence between successful parsing and invalid-argument protocol
handling. The CLI needs a maintained grammar library, but Quill's public
protocol must remain owned by Quill rather than the library.

## Decision

`github.com/alecthomas/kong` is pinned at v1.16.0 and is contained entirely in
`internal/cli`. Each command is one unexported grammar and operation struct;
the root command model declares every supported command in deterministic order.
A fresh model and parser are created for each `Runner.Run` invocation.

Kong owns only command grammar: flag and positional parsing, defaults, static
enum validation, and selected-command discovery. Quill owns help rendering,
help aliases, repository and init-target preparation, command execution,
stdout/stderr discipline, JSON envelopes, stable error categories, and exit
status mapping. Quill never calls Kong process-exiting helpers.

Documented flags use canonical long spellings such as `--format`. The prior
undocumented single-dash long forms are removed. Flags after `list` and
`explain` positional arguments are valid when Kong accepts them.

Architecture tests forbid Kong imports outside `internal/cli`.

## Consequences

- Every command's grammar is adjacent to its operation implementation.
- Parsing occurs once per invocation; invalid JSON-mode requests retain the
  format established by Kong's parse context.
- Help and errors remain deterministic Quill protocol output rather than
  terminal-width-dependent framework output.
- Kong is an explicit implementation dependency, not a public API or a
  dependency of engine, domain, report, profile, or adapter packages.

## Alternatives considered

### Keep the custom `flag` framework

Rejected. It required duplicated grammar, a mutable registry, and a second
parse for protocol recovery. Continued local maintenance would not improve
Quill's product behaviour.

### Adopt Cobra

Rejected. Cobra would solve command grammar but brings a broader command
framework whose normal execution, help, and error conventions Quill would need
to suppress. Kong maps directly to the small invocation-local grammar model
without command registration or process lifecycle ownership.

### Let Kong render help and exit

Rejected. That would make a dependency define Quill's public help, stream, and
exit behaviour. Quill must preserve its own deterministic protocol and cleanup
lifecycle.
