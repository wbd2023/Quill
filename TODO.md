# Quill roadmap

Quill's stable integration surface is the `quill` command and the repository-owned `STYLE.md`,
`quill.toml`, and `quill.lock` files. Keep implementation packages internal until a concrete
external Go consumer justifies a public package contract.

## Tool lifecycle

- Introduce version constraints and multiple pinned instances only when a real Pack needs a
  different version of the same tool. The current one-pin-per-tool invariant belongs in
  `policy.PinnedTools`.
- Replace the sealed `toolchain.InstallMethod` set with explicit registration only when a fifth
  installation ecosystem requires it. Keep the current typed methods while the set remains closed.
- Define a supported release platform matrix, archive format, and checksum publication contract
  before adding GoReleaser or distributing prebuilt Quill binaries.

## Pack composition

- Decide whether enabling a Pack should activate its default Rules. Profiles currently enable Packs
  and bind every active Rule separately. Any defaulting model must keep consumer-owned enforcement,
  scope, and requirement bindings explicit.
- Co-locate each shipped Pack's rule definitions and runtime binding registration. Preserve explicit
  composition: no `init` registration and no global mutable registry. Add a coverage test proving
  that every execution ID declared by a shipped Rule has one binding before changing ownership.
- Separate tools inspected by `ToolchainExecution` from tools required to run an execution. The
  generic preflight currently blocks a toolchain Rule before its Driver can report per-tool invalid
  diagnostics. Choose one reporting owner and remove the unreachable path.

## Profile and validation

- Validate every Rule `requirement_ids` entry against the loaded `STYLE.md` requirement set during
  Profile preparation. Syntax validation alone does not catch a well-formed but unknown ID.
- Consider extending repository exclusions from directory basenames to repository-relative file
  and path patterns. Keep exclusion at the file-walk trust boundary rather than duplicating it in
  Checks.
- Revisit the flat `[[rules]]` list only when its size creates demonstrated authoring errors. Any
  new shape must make duplicate bindings and Pack defaults easier to understand, not merely move
  TOML.

## Internal maintainability

- Consolidate repeated raw-TOML decoding primitives from the Pack policy packages into one small
  internal package. Keep Pack schemas, defaults, and validation with their owning Pack; do not build
  a base-Pack framework.
- Audit the shipped vocabulary Pack when a real shorthand gap appears. Preserve idiomatic Go forms
  such as `ctx`, `err`, `req`, and `db`.
- Consider one generic registry for the structurally identical Driver binding maps only if call
  sites become simpler. Four boring typed registries are preferable to noisy generic signatures.
