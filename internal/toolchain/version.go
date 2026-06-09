package toolchain

import (
	"fmt"

	"ciphera/tools/internal/style"
)

const (
	toolVersionGoCommand  VersionKind = "go_command"
	toolVersionBuildInfo  VersionKind = "build_info"
	toolVersionShellcheck VersionKind = "shellcheck"
	toolVersionNodeCLI    VersionKind = "node_cli"
)

type versionHandler func(
	runner CommandRunner,
	tool style.Tool,
	capability Capability,
	path string,
	environment map[string]string,
) (string, error)

func SupportsVersionKind(kind VersionKind) (supported bool) {
	_, supported = versionHandlers()[kind]
	return supported
}

func detectVersion(
	runner CommandRunner,
	tool style.Tool,
	capability Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	handler, found := versionHandlers()[capability.VersionKind]
	if !found {
		return "", fmt.Errorf("unsupported version detector %s", capability.VersionKind)
	}

	return handler(runner, tool, capability, path, environment)
}

func versionHandlers() (handlers map[VersionKind]versionHandler) {
	return map[VersionKind]versionHandler{
		toolVersionGoCommand:  goCommandVersion,
		toolVersionBuildInfo:  buildInfoVersion,
		toolVersionShellcheck: shellcheckVersion,
		toolVersionNodeCLI:    nodeCLIVersion,
	}
}

func goCommandVersion(
	runner CommandRunner,
	_ style.Tool,
	_ Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	return detectGoVersion(runner, path, environment)
}

func buildInfoVersion(
	_ CommandRunner,
	_ style.Tool,
	capability Capability,
	path string,
	_ map[string]string,
) (version string, err error) {
	return detectBuildInfoVersion(capability, path)
}

func shellcheckVersion(
	runner CommandRunner,
	_ style.Tool,
	_ Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	return detectCommandVersion(runner, path, "--version", environment, parseShellcheckVersion)
}

func nodeCLIVersion(
	runner CommandRunner,
	_ style.Tool,
	_ Capability,
	path string,
	environment map[string]string,
) (version string, err error) {
	return detectCommandVersion(runner, path, "--version", environment, parseSingleTokenVersion)
}
