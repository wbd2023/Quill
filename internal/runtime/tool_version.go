package runtime

import (
	"debug/buildinfo"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

const (
	toolVersionGoCommand  toolchain.VersionKind = "go_command"
	toolVersionBuildInfo  toolchain.VersionKind = "build_info"
	toolVersionShellcheck toolchain.VersionKind = "shellcheck"
	toolVersionNodeCLI    toolchain.VersionKind = "node_cli"
)

/* -------------------------------------- Version Detection ------------------------------------- */

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

type versionHandler func(
	tool contract.Tool,
	capability toolchain.Capability,
	path string,
	environment map[string]string,
) (string, error)

func versionHandlers() (handlers map[toolchain.VersionKind]versionHandler) {
	return map[toolchain.VersionKind]versionHandler{
		toolVersionGoCommand: func(
			_ contract.Tool,
			_ toolchain.Capability,
			path string,
			environment map[string]string,
		) (string, error) {
			return detectGoVersion(path, environment)
		},
		toolVersionBuildInfo: func(
			_ contract.Tool,
			capability toolchain.Capability,
			path string,
			_ map[string]string,
		) (string, error) {
			return detectBuildInfoVersion(capability, path)
		},
		toolVersionShellcheck: func(
			_ contract.Tool,
			_ toolchain.Capability,
			path string,
			environment map[string]string,
		) (string, error) {
			return detectCommandVersion(path, "--version", environment, parseShellcheckVersion)
		},
		toolVersionNodeCLI: func(
			_ contract.Tool,
			_ toolchain.Capability,
			path string,
			environment map[string]string,
		) (string, error) {
			return detectCommandVersion(path, "--version", environment, parseSingleTokenVersion)
		},
	}
}

func SupportsVersionKind(kind toolchain.VersionKind) (supported bool) {
	_, supported = versionHandlers()[kind]
	return supported
}

func detectGoVersion(
	commandPath string,
	environment map[string]string,
) (version string, err error) {
	result, err := RunCommand(CommandRequest{
		Directory:   ".",
		Environment: environment,
		Name:        commandPath,
		Arguments:   []string{"version"},
	})
	output, err := CommandOutput(result, err)
	if err != nil {
		return "", err
	}

	for _, field := range strings.Fields(output) {
		if !strings.HasPrefix(field, "go") {
			continue
		}

		version := normaliseVersion(strings.TrimPrefix(field, "go"))
		if version != "" {
			return version, nil
		}
	}

	return "", fmt.Errorf("could not parse go version")
}

func detectBuildInfoVersion(
	capability toolchain.Capability,
	path string,
) (version string, err error) {
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read embedded build info")
	}

	if capability.ModulePath != "" && info.Main.Path != capability.ModulePath {
		return "", fmt.Errorf("unexpected build target %s", info.Main.Path)
	}

	if info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "", fmt.Errorf("binary does not expose a pinned module version")
	}

	return info.Main.Version, nil
}

func detectCommandVersion(
	commandPath string,
	argument string,
	environment map[string]string,
	parse func(string) (string, error),
) (version string, err error) {
	result, err := RunCommand(CommandRequest{
		Directory:   ".",
		Environment: environment,
		Name:        commandPath,
		Arguments:   []string{argument},
	})
	output, err := CommandOutput(result, err)
	if err != nil {
		return "", err
	}

	return parse(output)
}

func inspectLocalToolVersion(
	tool contract.Tool,
	capability toolchain.Capability,
	path string,
) (version string, found bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, nil
		}

		return "", false, err
	}

	version, err = detectVersion(tool, capability, path, nil)
	if err != nil {
		return "", false, nil
	}

	return version, true, nil
}

func parseShellcheckVersion(output string) (version string, err error) {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "version:"); ok {
			return strings.TrimSpace(after), nil
		}
	}

	return "", fmt.Errorf("could not parse shellcheck version")
}

func parseSingleTokenVersion(output string) (version string, err error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return "", fmt.Errorf("could not parse version output")
	}

	return strings.TrimPrefix(fields[0], "v"), nil
}

/* ------------------------------------ Version Normalisation ----------------------------------- */

func matchesPinnedVersion(actual string, pinned string) (matches bool) {
	return normaliseVersion(actual) == normaliseVersion(pinned)
}

func normaliseVersion(version string) (normalised string) {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	end := len(version)
	for index, character := range version {
		if (character < '0' || character > '9') && character != '.' {
			end = index
			break
		}
	}

	version = version[:end]
	if version == "" {
		return ""
	}

	parts := strings.Split(version, ".")
	for _, piece := range parts {
		if _, err := strconv.Atoi(piece); err != nil {
			return ""
		}
	}

	return strings.Join(parts, ".")
}
