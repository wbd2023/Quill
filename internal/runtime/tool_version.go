package runtime

import (
	"debug/buildinfo"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ciphera/tools/internal/contract"
)

/* -------------------------------------- Version Detection ------------------------------------- */

func detectVersion(
	tool contract.Tool,
	path string,
	environment map[string]string,
) (version string, err error) {
	switch tool.VersionKind {
	case contract.ToolVersionGoCommand:
		return detectGoVersion(path, environment)

	case contract.ToolVersionBuildInfo:
		return detectBuildInfoVersion(tool, path)

	case contract.ToolVersionShellcheck:
		return detectCommandVersion(path, "--version", environment, parseShellcheckVersion)

	case contract.ToolVersionNodeCLI:
		return detectCommandVersion(path, "--version", environment, parseSingleTokenVersion)

	default:
		return "", fmt.Errorf("unsupported version detector %s", tool.VersionKind)
	}
}

func detectGoVersion(
	commandPath string,
	environment map[string]string,
) (version string, err error) {
	output, err := RunCommand(".", environment, commandPath, "version")
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

func detectBuildInfoVersion(tool contract.Tool, path string) (version string, err error) {
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read embedded build info")
	}

	if tool.ModulePath != "" && info.Main.Path != tool.ModulePath {
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
	output, err := RunCommand(".", environment, commandPath, argument)
	if err != nil {
		return "", err
	}

	return parse(output)
}

func inspectLocalToolVersion(
	tool contract.Tool,
	path string,
) (version string, found bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, nil
		}

		return "", false, err
	}

	version, err = detectVersion(tool, path, nil)
	if err != nil {
		return "", false, nil
	}

	return version, true, nil
}

func parseShellcheckVersion(output string) (version string, err error) {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "version:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "version:")), nil
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

func normaliseVersion(value string) (normalised string) {
	value = strings.TrimPrefix(strings.TrimSpace(value), "v")
	end := len(value)
	for index, runeValue := range value {
		if (runeValue < '0' || runeValue > '9') && runeValue != '.' {
			end = index
			break
		}
	}

	value = value[:end]
	if value == "" {
		return ""
	}

	parts := strings.Split(value, ".")
	for _, piece := range parts {
		if _, err := strconv.Atoi(piece); err != nil {
			return ""
		}
	}

	return strings.Join(parts, ".")
}
