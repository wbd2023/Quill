package toolchain

import (
	"debug/buildinfo"
	"fmt"
	"strconv"
	"strings"

	"ciphera/tools/internal/runtime"
)

/* ------------------------------------------ Detection ----------------------------------------- */

// detectVersion returns the installed version of the binary at path.
func detectVersion(
	environment map[string]string,
	path string,
	method VersionMethod,
) (version string, err error) {
	switch method := method.(type) {
	case GoVersion:
		return detectCommandVersion(environment, path, "version", parseGoVersion)

	case ModuleVersion:
		return detectModuleVersion(path, method)

	case PrefixedLineVersion:
		return detectCommandVersion(environment, path, "--version", parsePrefixedLineVersion)

	case FirstTokenVersion:
		return detectCommandVersion(environment, path, "--version", parseSingleTokenVersion)

	default:
		return "", fmt.Errorf("unsupported version method %T", method)
	}
}

// detectModuleVersion reads embedded build info; ModulePath, if set, must match the binary's main
// module.
func detectModuleVersion(path string, method ModuleVersion) (version string, err error) {
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read embedded build info")
	}

	if method.ModulePath != "" && info.Main.Path != method.ModulePath {
		return "", fmt.Errorf("unexpected build target %s", info.Main.Path)
	}

	if info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "", fmt.Errorf("binary does not expose a pinned module version")
	}

	return info.Main.Version, nil
}

// detectCommandVersion runs the binary at path with argument and parses the output.
func detectCommandVersion(
	environment map[string]string,
	path string,
	argument string,
	parse func(string) (string, error),
) (version string, err error) {
	result, err := runtime.RunCommand(runtime.CommandRequest{
		Environment: environment,
		Name:        path,
		Arguments:   []string{argument},
	})
	if err != nil {
		return "", err
	}

	return parse(result.Output)
}

/* ------------------------------------------- Parsing ------------------------------------------ */

// parseGoVersion extracts the goX.Y.Z token from `go version` output.
func parseGoVersion(output string) (version string, err error) {
	for field := range strings.FieldsSeq(output) {
		if !strings.HasPrefix(field, "go") {
			continue
		}

		if version = normaliseVersion(strings.TrimPrefix(field, "go")); version != "" {
			return version, nil
		}
	}

	return "", fmt.Errorf("could not parse go version")
}

// parsePrefixedLineVersion finds the first "version:" prefixed line and returns its value.
func parsePrefixedLineVersion(output string) (version string, err error) {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "version:"); ok {
			return strings.TrimSpace(after), nil
		}
	}

	return "", fmt.Errorf("could not parse version from prefixed line")
}

// parseSingleTokenVersion returns the first whitespace-delimited token, stripped of a leading v.
func parseSingleTokenVersion(output string) (version string, err error) {
	for field := range strings.FieldsSeq(output) {
		return strings.TrimPrefix(field, "v"), nil
	}

	return "", fmt.Errorf("could not parse version output")
}

/* ---------------------------------------- Normalisation --------------------------------------- */

// normaliseVersion strips a leading v and truncates at the first non-numeric, non-dot character,
// returning the dot-separated numeric prefix (e.g. "v1.2.3-rc1" becomes "1.2.3"). Returns the
// empty string if any segment is non-numeric.
func normaliseVersion(version string) (normalised string) {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	end := len(version)
	for index, char := range version {
		if (char < '0' || char > '9') && char != '.' {
			end = index
			break
		}
	}

	version = version[:end]
	if version == "" {
		return ""
	}

	parts := strings.Split(version, ".")
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return ""
		}
	}

	return strings.Join(parts, ".")
}
