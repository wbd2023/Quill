# Security

Quill installs and executes repository-pinned development tools. This document defines the trust
boundary for maintainers and consuming repositories.

## Reporting a vulnerability

Do not open a public issue for an undisclosed vulnerability. Use GitHub private vulnerability
reporting for `wbd2023/Quill` when available, or contact the maintainer privately through the
contact details on the maintainer's GitHub profile.

Include the affected commit or release, operating system, command, Profile and lockfile fragments,
and a minimal reproduction. Avoid including live credentials or unrelated repository contents.

The maintainer will acknowledge a complete report, assess impact, and coordinate disclosure after a
fix is available. No response-time guarantee is currently offered for this pre-1.0 project.

## Supported versions

Before the first tagged release, only the current `main` branch is supported. After releases begin,
security fixes will target the latest release unless a release announcement states otherwise.

## Trust model

Quill is not a sandbox. Running Quill grants its process the same filesystem and process permissions
as the invoking user.

Treat these inputs as trusted code or policy:

- the Quill executable;
- the consuming repository's `quill.toml` and `quill.lock`;
- shipped Pack tool capabilities;
- installed third-party tools;
- repository Make targets and scripts invoked around Quill.

Review changes to these inputs before running Quill on an untrusted checkout. Run CI jobs with the
minimum repository and token permissions required.

## Tool installation

Quill installs tools beneath a repository-local cache. Archive-installed tools must use HTTPS and
must match the selected platform hash in `quill.lock`.

The installer rejects:

- missing or incorrect checksums;
- archive paths that escape the destination;
- symbolic links and hard links;
- unsupported archive entries;
- downloads or extracted content above configured size limits;
- unsafe replacement of an existing installation.

`quill lock` writes a complete replacement atomically with mode `0644`. Review lockfile changes in
the same way as dependency lockfile changes. Never accept a new hash solely to silence a failed
verification.

## Process execution

Checks and fixes invoke pinned external tools with explicit working directories, environments,
timeouts, and output limits. Tool output is untrusted data. Quill distinguishes checker findings
from command-execution failures, but it cannot prevent a malicious executable from acting with the
invoking user's permissions.

Use reviewed, pinned tool versions. Do not substitute ambient executables in CI. Do not use
`@latest` or an unverified release binary for a protected branch.

## Filesystem safety

Quill discovers a repository by locating both `STYLE.md` and `quill.toml`. Use `--repo-root` when
running from automation that could otherwise discover the wrong parent repository.

Generated files, dependency trees, caches, and fixture directories should be excluded in the
Profile. Exclusions reduce unintended scanning and rewriting; they do not create a security
sandbox.

## Secrets and privacy

Quill scans source and may include file paths, line numbers, checker output, and source-derived
messages in text or JSON reports. Treat reports and CI logs according to the sensitivity of the
checked repository.

Do not place secrets in `quill.toml`, `quill.lock`, command output fixtures, or bug reports. Quill's
secret checks are defence-in-depth and do not replace a dedicated secret-management process.

## Security regression requirements

Changes to downloads, hashes, archive extraction, process execution, repository discovery, atomic
writes, or permission handling require tests for the affected rejection and success paths. A
security check must fail closed when integrity or path safety cannot be established.
