# Quill

Quill turns a repository's human-authored `STYLE.md` and machine-readable `quill.toml` Profile into
executable style checks. It resolves Pack defaults, validates Rule bindings, installs pinned tools,
runs checks and safe fixes, writes `quill.lock`, and reports requirement coverage.

Quill is a CLI application. Its Go packages are internal implementation details.
The stable integration surface is the `quill` command, its repository formats,
and the local External Pack manifest and subprocess protocol.

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
go install github.com/wbd2023/quill/cmd/quill@v0.1.0
```

Do not use `@latest` in CI. Pin a reviewed module release.

## Integrate

Quill is a CLI-first, language-neutral product. Pin a reviewed `quill`
release, invoke it for a local repository, and consume its documented command
output and exit status. This works from any language or automation environment
that can execute a supported Quill artefact.

For a Go repository, `go tool quill` is an executable integration mechanism,
not an importable Go API. See [docs/product.md](docs/product.md) for the
supported interfaces and planned machine protocol.

## Versioning and releases

`quill version` reads the main module version from `runtime/debug.BuildInfo`. Go 1.24 stamps that
version from the repository tag or commit and appends `+dirty` when the checkout has uncommitted
changes. Builds without usable module or VCS version information report `(devel)`.

The Git tag is the only version source. Quill has no separate version file and
does not inject a second version with linker flags. This keeps local tagged
builds, release archives, and Go 1.24 `tool` dependencies on the same version
contract.

For each release:

1. Choose and review the semantic version.
2. From a clean release commit, run `make lint`, `make test`, and `go vet ./...`.
3. Create the matching semantic-version tag locally.
4. Run `GOFLAGS=-buildvcs=true make build` and confirm that `./bin/quill version`
   prints the tag.
5. Push the tag. Release CI builds, verifies, checksums, and publishes the
   supported binary archives.
6. Download the intended archive, verify its SHA-256 entry in the published
   checksum file, and confirm `quill version` prints the tag.

Release CI publishes these archives:

- Linux amd64 and arm64 (`.tar.gz`);
- macOS amd64 and arm64 (`.tar.gz`); and
- Windows amd64 (`.zip`).

Each archive contains only `quill` or `quill.exe`. The release also contains
`quill_<version>_checksums.txt`, with one SHA-256 digest per archive.

## Repository contract

A checked repository owns the policy that Quill executes:

- `STYLE.md` is the human source of truth. Requirements carry stable IDs in hidden
  `<!-- style: id=... -->` metadata.
- `quill.toml` is the executable Profile. It selects Packs, binds Rules to requirements and scopes,
  declares Targets and file sets, and pins tool versions.
- `quill.lock` records verified per-platform hashes for archive-installed tools.

Quill has no built-in knowledge of a consuming repository's package layout, scope names, naming
vocabulary, or policy values. Those decisions belong in that repository's Profile.

When `--repository-root` is omitted, Quill walks upward from the current
directory until it finds both `STYLE.md` and `quill.toml`.

## External Packs

A repository can add a local external Pack through `[[pack_sources]]` in its
Profile. The Pack's `pack.toml` manifest and `quill-pack-v1` subprocess protocol
are supported first-release interfaces. See
[docs/pack-protocol.md](docs/pack-protocol.md) before authoring or reviewing one.
External Pack executables run with the invoking user's permissions; Quill does
not sandbox them.

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
- `quill init` creates a minimal `STYLE.md` and `quill.toml` in an empty directory.
- `quill list <packs|rules|tools|scopes>` lists compiled Profile metadata.
- `quill explain <rule:ID>` explains an active Rule.

Use `quill help <command>` for command-specific flags. Canonical flags use
the long `--flag` spelling.

Common examples:

```sh
quill check --repository-root . --mode required
quill check --repository-root . --mode all --strict-recommendations --verbose
quill check --repository-root . --scope all --format json
quill fix --repository-root . --scope all --format json
quill doctor --repository-root . --format json
quill coverage --repository-root . --format json
quill install --repository-root . --format json
quill lock --repository-root . --format json
quill version
quill init --repository-root .
quill list packs --format json --repository-root .
quill explain rule:profile/enforcement-levels --format json --repository-root .
```

Exit codes:

- `0`: command completed successfully and no selected failure requires a non-zero result.
- `1`: a selected Rule failed, a Rule errored, or command execution failed.
- `2`: command-line usage was invalid.

JSON output is intended for automation. Text output is intended for people.
See [docs/cli-protocol.md](docs/cli-protocol.md) for stream, envelope, error,
cancellation, and compatibility rules.

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

Quill is a modular monolith with one CLI entrypoint and private implementation
packages. Its public product boundary is the CLI and repository formats. See
[docs/README.md](docs/README.md) for the documentation map,
[docs/mvp.md](docs/mvp.md) for the ideal first-release target,
[docs/architecture.md](docs/architecture.md) for package ownership and runtime
flow, [docs/product.md](docs/product.md) for supported interfaces, and
[ADR 0004](docs/adr/0004-cli-first-language-neutral-product.md) for the
accepted boundary decision.

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

Quill has one command and private implementation packages under `internal/`.
ADR 0004 is complete: no production Go package exists at the repository root.

- `cmd/quill` contains only process entrypoint wiring.
- `internal/cli` owns argument parsing, output streams, exit codes, and the
  language-neutral command protocol.
- `internal/engine` is the private application facade for repository workflows,
  including check, fix, inspection, coverage, installation, locking,
  metadata, explanation, and initialization.
- `internal/profile` owns the Profile model and persisted codec, then loads,
  validates, and compiles Profiles into executable plans.
- `internal/pack` defines Pack contracts and resolution.
- `internal/pack/shipped` assembles Shipped Pack declarations and Pack-local
  runtime bindings.
- `internal/checks` contains concrete repository Checks.
- `internal/execution` runs compiled Jobs through Drivers.
- `internal/installer`, `internal/toolchain`, and `internal/process` own
  external-tool boundaries.
- `internal/styleguide` parses STYLE.md requirement metadata.
- `internal/coverage` derives requirement coverage from the compiled Profile.
- `internal/report` renders command results as CLI text and JSON.

Architecture tests under `internal/architecture` enforce important import and
ownership boundaries.

## Security

Quill executes repository policy, local external Packs, and downloaded tools on
the host. It is not a sandbox. Review `quill.toml`, `quill.lock`, local Pack
sources, and changes to Shipped Pack capabilities before running Quill on an
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
