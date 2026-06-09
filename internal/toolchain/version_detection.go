package toolchain

import (
	"debug/buildinfo"
	"fmt"
	"strings"
)

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
	capability Capability,
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
