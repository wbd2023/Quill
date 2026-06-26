package toolchain

import "fmt"

// SupportsVersionKind reports whether kind names a known version-detection strategy.
func SupportsVersionKind(kind VersionKind) (supported bool) {
	switch kind {

	case VersionKindGoCommand,
		VersionKindBuildInfo,
		VersionKindShellcheck,
		VersionKindNodeCLI:
		return true
	}

	return false
}

func detectVersion(
	runner CommandRunner,
	capability Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	switch capability.VersionKind {

	case VersionKindGoCommand:
		return detectGoVersion(runner, path, environment)

	case VersionKindBuildInfo:
		return detectBuildInfoVersion(capability, path)

	case VersionKindShellcheck:
		return detectCommandVersion(runner, path, "--version", environment, parseShellcheckVersion)

	case VersionKindNodeCLI:
		return detectCommandVersion(runner, path, "--version", environment, parseSingleTokenVersion)

	default:
		return "", fmt.Errorf("unsupported version detector %q", capability.VersionKind)
	}
}
