# Quill rewrite references

This directory currently retains only language-independent product decisions and final-Go
contract evidence needed to design the Rust rewrite. These documents are inputs to the redesign,
not automatic specifications for the new implementation.

## Product and interface evidence

- [product.md](product.md) records the current product direction and supported integration model.
- [cli-protocol.md](cli-protocol.md) records the final-Go machine command interface.
- [pack-protocol.md](pack-protocol.md) records the final-Go local External Pack interface.

The protocol documents were added after the repository's `v0.1.0` tag. Before assigning
compatibility status or reusing version number 1, determine whether these post-tag interfaces were
published or adopted elsewhere.

## Retained decisions

- [ADR 0004](adr/0004-cli-first-language-neutral-product.md) records the CLI-first,
  language-neutral product boundary.
- [ADR 0006](adr/0006-external-pack-protocol-is-a-public-interface.md) records the subprocess-based
  External Pack extension boundary.

Both decisions must be revalidated against the clean-slate product model before they constrain Rust
module structure or protocol compatibility.

## Historical implementation

The final Go implementation and its architecture, review programme, delivery roadmap, tests, and
superseded decisions remain available at parent commit:

```text
3ed482e569b92cd6b4b7f1be5a0b80d64fbaa4e5
```

Do not restore those files into the active tree merely for convenient browsing. Inspect them through
Git when a product requirement or difficult edge case needs evidence.
