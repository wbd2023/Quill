package toolchain

import (
	"debug/buildinfo"
	"fmt"
	"strings"
)

/* ------------------------------------------ Dispatch ------------------------------------------ */

// detectVersion dispatches to the version-detection strategy carried by spec.
func detectVersion(
	runner CommandRunner,
	spec VersionSpec,
	path string,
	environment map[string]string,
) (version string, err error) {
	switch versionSpec := spec.(type) {

	case GoCommandVersion:
		return detectGoVersion(runner, path, environment)

	case BuildInfoVersion:
		return detectBuildInfoVersion(versionSpec, path)

	case PrefixedLineVersion:
		return detectCommandVersion(
			runner, path, "--version", environment, parsePrefixedLineVersion,
		)

	case FirstTokenVersion:
		return detectCommandVersion(runner, path, "--version", environment, parseSingleTokenVersion)

	default:
		return "", fmt.Errorf("unsupported version spec %T", spec)
	}
}

/* -------------------------------------- Command Versions -------------------------------------- */

func detectGoVersion(
	runner CommandRunner,
	commandPath string,
	environment map[string]string,
) (version string, err error) {
	if runner == nil {
		return "", fmt.Errorf("version command runner is not configured")
	}

	output, err := runner(CommandRequest{
		Environment: environment,
		Name:        commandPath,
		Arguments:   []string{"version"},
	})
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

/* ----------------------------------------- Build Info ----------------------------------------- */

func detectBuildInfoVersion(
	spec BuildInfoVersion,
	path string,
) (version string, err error) {
	info, err := buildinfo.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read embedded build info")
	}

	if spec.ModulePath != "" && info.Main.Path != spec.ModulePath {
		return "", fmt.Errorf("unexpected build target %s", info.Main.Path)
	}

	if info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "", fmt.Errorf("binary does not expose a pinned module version")
	}

	return info.Main.Version, nil
}

/* --------------------------------------- Shared Helpers --------------------------------------- */

func detectCommandVersion(
	runner CommandRunner,
	commandPath string,
	argument string,
	environment map[string]string,
	parse func(string) (string, error),
) (version string, err error) {
	if runner == nil {
		return "", fmt.Errorf("version command runner is not configured")
	}

	output, err := runner(CommandRequest{
		Environment: environment,
		Name:        commandPath,
		Arguments:   []string{argument},
	})
	if err != nil {
		return "", err
	}

	return parse(output)
}
