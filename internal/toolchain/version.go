package toolchain

import (
	"context"
	"debug/buildinfo"
	"fmt"
	"strconv"
	"strings"
)

/* ------------------------------------------ Detection ----------------------------------------- */

// DetectByCommand returns a VersionMethod that runs the tool with argument and extracts the version
// from its output using extract.
func DetectByCommand(argument string, extract func(string) (string, error)) (method VersionMethod) {
	return func(
		ctx context.Context,
		runner CommandRunner,
		environment map[string]string,
		path string,
	) (version string, err error) {
		output, err := runner.Run(ctx, environment, path, []string{argument})
		if err != nil {
			return "", err
		}

		return extract(output)
	}
}

// DetectByGoBinary returns a VersionMethod that reads the version embedded in a Go binary's build
// info. ModulePath, if set, must match the binary's main module.
func DetectByGoBinary(modulePath string) (method VersionMethod) {
	return func(
		ctx context.Context,
		runner CommandRunner,
		environment map[string]string,
		path string,
	) (version string, err error) {
		info, err := buildinfo.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("could not read embedded build info")
		}

		if modulePath != "" && info.Main.Path != modulePath {
			return "", fmt.Errorf("unexpected build target %s", info.Main.Path)
		}

		if info.Main.Version == "" || info.Main.Version == "(devel)" {
			return "", fmt.Errorf("binary does not expose a pinned module version")
		}

		return info.Main.Version, nil
	}
}

/* ------------------------------------------- Parsing ------------------------------------------ */

// ExtractGoToken extracts the goX.Y.Z token from `go version` output.
func ExtractGoToken(output string) (version string, err error) {
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

// ExtractPrefixedLine finds the first "version:" prefixed line and returns its value.
func ExtractPrefixedLine(output string) (version string, err error) {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "version:"); ok {
			return strings.TrimSpace(after), nil
		}
	}

	return "", fmt.Errorf("could not parse version from prefixed line")
}

// ExtractFirstToken returns the first whitespace-delimited token, stripped of a leading v.
func ExtractFirstToken(output string) (version string, err error) {
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
