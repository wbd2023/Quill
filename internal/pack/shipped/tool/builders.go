package tool

import "ciphera/tools/internal/toolchain"

func buildBuiltin(
	id string,
	name string,
	command string,
	version toolchain.VersionSpec,
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
		Version: toolchain.BuildInfoVersion{ModulePath: modulePath},
		Install: toolchain.GoBinaryInstall{Source: installSource},
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
		Version: toolchain.FirstTokenVersion{},
		Install: toolchain.NodePackageInstall{Source: installSource},
	}
}

func buildShellcheckArchive() (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:      Shellcheck,
		Name:    "shellcheck",
		Command: "shellcheck",
		Version: toolchain.PrefixedLineVersion{},
		Install: toolchain.ArchiveInstall{
			Spec: toolchain.ArchiveSpec{
				URLFormat: "https://github.com/koalaman/shellcheck/releases/download/" +
					"v%[1]s/shellcheck-v%[1]s.%[2]s.tar.xz",
				BinaryPathFormat: "shellcheck-v%[1]s/shellcheck",
				Platforms: map[string]string{
					"darwin/amd64": "darwin.x86_64",
					"darwin/arm64": "darwin.aarch64",
					"linux/amd64":  "linux.x86_64",
					"linux/arm64":  "linux.aarch64",
				},
			},
		},
	}
}
