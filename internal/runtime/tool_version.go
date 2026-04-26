package runtime

import (
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

const (
	toolVersionGoCommand  toolchain.VersionKind = "go_command"
	toolVersionBuildInfo  toolchain.VersionKind = "build_info"
	toolVersionShellcheck toolchain.VersionKind = "shellcheck"
	toolVersionNodeCLI    toolchain.VersionKind = "node_cli"
)

type versionHandler func(
	tool contract.Tool,
	capability toolchain.Capability,
	path string,
	environment map[string]string,
) (string, error)

func SupportsVersionKind(kind toolchain.VersionKind) (supported bool) {
	_, supported = versionHandlers()[kind]
	return supported
}

func detectVersion(
	tool contract.Tool,
	capability toolchain.Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	handler, found := versionHandlers()[capability.VersionKind]
	if !found {
		return "", fmt.Errorf("unsupported version detector %s", capability.VersionKind)
	}

	return handler(tool, capability, path, environment)
}

func versionHandlers() (handlers map[toolchain.VersionKind]versionHandler) {
	return map[toolchain.VersionKind]versionHandler{
		toolVersionGoCommand:  goCommandVersion,
		toolVersionBuildInfo:  buildInfoVersion,
		toolVersionShellcheck: shellcheckVersion,
		toolVersionNodeCLI:    nodeCLIVersion,
	}
}

func goCommandVersion(
	_ contract.Tool,
	_ toolchain.Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	return detectGoVersion(path, environment)
}

func buildInfoVersion(
	_ contract.Tool,
	capability toolchain.Capability,
	path string,
	_ map[string]string,
) (version string, err error) {
	return detectBuildInfoVersion(capability, path)
}

func shellcheckVersion(
	_ contract.Tool,
	_ toolchain.Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	return detectCommandVersion(path, "--version", environment, parseShellcheckVersion)
}

func nodeCLIVersion(
	_ contract.Tool,
	_ toolchain.Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	return detectCommandVersion(path, "--version", environment, parseSingleTokenVersion)
}
