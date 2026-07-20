# Quill

Quill turns a repository's human-authored `STYLE.md` and machine-readable `quill.toml` Profile into
executable style checks. It resolves Pack defaults, validates Rule bindings, installs pinned tools,
runs checks and safe fixes, writes `quill.lock`, and reports requirement coverage.

Quill is a CLI application. Its Go packages are internal implementation details; the stable
integration surface is the `quill` command and the repository-owned files it consumes.

## Status

Quill is pre-1.0. The CLI and file formats are being prepared for their first standalone release.
Until a release is tagged, build from a reviewed commit rather than relying on an unpinned branch.

## Requirements

- Go 1.24.5, matching the self-check Profile's exact toolchain pin
- Node.js 20 when installing or running Node-based tools
- A POSIX environment for the repository Make targets

## Install

Build the current checkout:

```sh
make build
./bin/quill version
```

Install the current release:

```sh
go install github.com/wbd2023/Quill/cmd/quill@v0.1.0
```

Do not use `@latest` in CI. Pin a reviewed module release.

## Versioning and releases

`quill version` reads the main module version from `runtime/debug.BuildInfo`. Go 1.24 stamps that
version from the repository tag or commit and appends `+dirty` when the checkout has uncommitted
changes. Builds without usable module or VCS version information report `devel`.

The Git tag is the only version source. Quill has no separate version file and does not inject a
second version with linker flags. This keeps local tagged builds, `go install ...@vX.Y.Z`, and Go
1.24 `tool` dependencies on the same version contract.

For each release:

1. Choose and review the semantic version.
2. From a clean release commit, run `make lint`, `make test`, and `go vet ./...`.
3. Create the matching semantic-version tag locally.
4. Run `GOFLAGS=-buildvcs=true make build` and confirm that `./bin/quill version` prints the tag.
5. Push the tag. Tag CI repeats the VCS-enabled build and version assertion.
6. Verify the pinned `go install` command, then publish release notes.

Releases publish source only. Quill does not publish archive binaries until their supported platform
matrix, checksum generation, and installation contract are defined.

## Repository contract

A checked repository owns the policy that Quill executes:

- `STYLE.md` is the human source of truth. Requirements carry stable IDs in hidden
  `<!-- style: id=... -->` metadata.
- `quill.toml` is the executable Profile. It selects Packs, binds Rules to requirements and scopes,
  declares Targets and file sets, and pins tool versions.
- `quill.lock` records verified per-platform hashes for archive-installed tools.

Quill has no built-in knowledge of a consuming repository's package layout, scope names, naming
vocabulary, or policy values. Those decisions belong in that repository's Profile.

When `--repo-root` is omitted, Quill walks upward from the current directory until it finds both
`STYLE.md` and `quill.toml`.

## CLI

```text
quill <command> [flags]
```

Commands:

- `quill check` runs selected Rules.
- `quill fix` runs safe fixes for selected Rules.
- `quill doctor` inspects pinned tools and reports missing or wrong versions.
- `quill coverage` maps STYLE.md requirements to automated, review-only, and deferred coverage.
- `quill install` installs or refreshes pinned tools in the repository-local cache.
- `quill lock` resolves archive-tool hashes and atomically rewrites `quill.lock`.
- `quill version` prints the version recorded by the Go toolchain.

Use `quill help <command>` for command-specific flags.

Common examples:

```sh
quill check --mode required
quill check --mode all --strict-recommendations --verbose
quill check --scope all --format json
quill fix --scope all
quill doctor --format json
quill coverage --format json
quill lock
```

Exit codes:

- `0`: command completed successfully and no selected failure requires a non-zero result.
- `1`: a selected Rule failed, a Rule errored, or command execution failed.
- `2`: command-line usage was invalid.

JSON output is intended for automation. Text output is intended for people.

## Profile model

A Profile contains seven main areas:

1. Repository roots, scopes, exclusions, and generated-file markers.
2. The STYLE.md path and path-role classifications.
3. Named file sets for repository scans.
4. Language Targets and their working directories.
5. Pinned tool versions and execution limits.
6. Enabled Packs and Pack Policy values.
7. Rule bindings, enforcement levels, scopes, and requirement IDs.

Packs provide reusable checker capabilities and defaults. The Profile decides which capabilities
are active. Drivers execute resolved jobs; Checks implement repository-specific observations
without owning Profile policy.

The repository's own `quill.toml` and `STYLE.md` are a complete self-checking example.

## Architecture

Quill is a modular monolith with one CLI entrypoint and private implementation packages. See
[docs/architecture.md](docs/architecture.md) for package ownership and runtime flow. The package
composition and public interface decisions are recorded in
[ADR 0001](docs/adr/0001-separate-shipped-declarations-from-check-execution.md) and
[ADR 0002](docs/adr/0002-cli-and-files-are-the-public-interface.md).

## Development

Install the pinned development tools once:

```sh
make style-install
```

Run the required gate:

```sh
make lint-required
```

Run the complete strict gate and tests:

```sh
make lint
make test
```

Build and smoke-test the command:

```sh
go build -o /tmp/quill ./cmd/quill
/tmp/quill help
/tmp/quill help check
/tmp/quill help lock
```

The repository keeps build products and installed tools under ignored repository-local
directories. Development commands do not mutate the global GOPATH tool directory.

## Package map

Quill is a modular monolith with one command and private packages under `internal/`.

- `cmd/quill` contains only process entrypoint wiring.
- `internal/cli` owns argument parsing, output streams, and exit codes.
- `internal/engine` is the operation facade for check, fix, doctor, coverage, install, and lock.
- `internal/profile` loads, validates, and compiles Profiles.
- `internal/profile/toml` owns the persisted Profile codec.
- `internal/pack` defines Pack contracts and resolution.
- `internal/pack/shipped` assembles built-in Pack declarations.
- `internal/checks` contains concrete repository Checks and Pack Policy codecs.
- `internal/execution` runs resolved jobs through Drivers.
- `internal/installer`, `internal/toolchain`, and `internal/process` own external-tool boundaries.
- `internal/styleguide` parses STYLE.md requirement metadata.
- `internal/coverage` derives requirement coverage from the compiled Profile.
- `internal/report` renders text and JSON views.

Architecture tests under `internal/architecture` enforce the intended import direction.

## Security

Quill executes repository policy and downloaded tools on the host. It is not a sandbox. Review
`quill.toml`, `quill.lock`, and changes to shipped tool capabilities before running Quill on an
untrusted checkout.

Archive downloads are HTTPS-only and are verified against lockfile hashes. Installer tests defend
archive traversal, links, oversized downloads, and checksum mismatches. See
[SECURITY.md](SECURITY.md) for the trust model and vulnerability reporting process.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Contributions require recorded acceptance of the
[Contributor Licence Agreement](CLA.md) and must preserve CLI contracts, Profile validation,
package boundaries, and installer security properties.

## Licence

Quill is licensed under the [Apache License, Version 2.0](LICENSE). Its SPDX licence identifier is
`Apache-2.0`. See [NOTICE](NOTICE) for creator attribution.
