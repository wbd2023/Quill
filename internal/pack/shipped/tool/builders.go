package tool

import "github.com/wbd2023/Quill/internal/toolchain"

func buildBuiltin(
	id string,
	name string,
	command string,
	version toolchain.VersionMethod,
) (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:      id,
		Name:    name,
		Command: command,
		Version: version,
		Install: toolchain.NoInstall{},
	}
}

func buildGoBinary(
	id string,
	name string,
	command string,
	modulePath string,
	installSource string,
) (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:      id,
		Name:    name,
		Command: command,
		Version: toolchain.DetectByGoBinary(modulePath),
		Install: toolchain.GoInstall{Source: installSource},
	}
}

func buildNodePackage(
	id string,
	name string,
	command string,
	installSource string,
) (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:      id,
		Name:    name,
		Command: command,
		Version: toolchain.DetectByCommand("--version", toolchain.ExtractFirstToken),
		Install: toolchain.NpmInstall{Source: installSource},
	}
}

func buildShellcheckArchive() (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:      Shellcheck,
		Name:    "shellcheck",
		Command: "shellcheck",
		Version: toolchain.DetectByCommand("--version", toolchain.ExtractPrefixedLine),
		Install: toolchain.GitHubInstall{
			Owner:      "koalaman",
			Repository: "shellcheck",
			Tag:        "v%s",
			Asset:      "shellcheck-%s.%s.tar.xz",
			Path:       "shellcheck-%s/shellcheck",
			Platforms: map[string]string{
				"darwin/amd64": "darwin.x86_64",
				"darwin/arm64": "darwin.aarch64",
				"linux/amd64":  "linux.x86_64",
				"linux/arm64":  "linux.aarch64",
			},
		},
	}
}
