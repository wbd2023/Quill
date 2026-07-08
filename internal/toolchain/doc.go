// Package toolchain inspects and describes pinned external tools. It resolves a tool's executable
// path, detects its installed version, and reports whether that version matches the profile pin.
// Capabilities carry the per-tool metadata (version-detection strategy, install strategy, archive
// spec) consumed by the inspector and installer.
package toolchain
