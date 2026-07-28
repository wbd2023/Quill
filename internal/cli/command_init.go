package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

/* ---------------------------------------- Init Command ---------------------------------------- */

const (
	initSummary     = "create a minimal STYLE.md and quill.toml"
	defaultPreset   = "minimal"
	styleFileName   = "STYLE.md"
	profileFileName = "quill.toml"
)

const (
	configDirectoryPermission os.FileMode = 0o755
	configFilePermission      os.FileMode = 0o644
)

func runInit(_ context.Context, tool Tool, options initOptions) (exitCode int) {
	preset, ok := initPreset(options.preset)
	if !ok {
		tool.writeError(fmt.Errorf("unsupported preset %q", options.preset))
		return 1
	}

	stylePath := filepath.Join(options.repoRoot, styleFileName)
	profilePath := filepath.Join(options.repoRoot, profileFileName)

	occupied, err := policyFilesOccupied(stylePath, profilePath)
	if err != nil {
		tool.writeError(fmt.Errorf("check target directory %q: %w", options.repoRoot, err))
		return 1
	}
	if occupied {
		tool.writeError(fmt.Errorf(
			"refusing to overwrite existing policy files in %q: "+
				"STYLE.md and quill.toml already present",
			options.repoRoot,
		))
		return 1
	}

	if err := os.MkdirAll(options.repoRoot, configDirectoryPermission); err != nil {
		tool.writeError(fmt.Errorf("create target directory %q: %w", options.repoRoot, err))
		return 1
	}
	if err := writePolicyFiles(
		stylePath,
		profilePath,
		preset.styleGuide,
		preset.profile,
	); err != nil {
		tool.writeError(fmt.Errorf("write policy files in %q: %w", options.repoRoot, err))
		return 1
	}

	if _, err := fmt.Fprintf(
		tool.stdout,
		"Initialised Quill in %s (preset %s)\n",
		options.repoRoot,
		options.preset,
	); err != nil {
		tool.writeError(fmt.Errorf("write initialisation summary: %w", err))
		return 1
	}
	if _, err := fmt.Fprintln(tool.stdout, "  wrote STYLE.md, quill.toml"); err != nil {
		tool.writeError(fmt.Errorf("write initialisation files: %w", err))
		return 1
	}
	return 0
}

// policyFilesOccupied reports whether any policy file path already names an entry. It uses Lstat
// so it does not follow symlinks: a dangling or valid symlink at a policy path counts as occupied
// and is refused, preventing init from writing through an attacker-controlled link.
func policyFilesOccupied(paths ...string) (occupied bool, err error) {
	for _, path := range paths {
		switch _, statErr := os.Lstat(path); {
		case statErr == nil:
			return true, nil
		case errors.Is(statErr, fs.ErrNotExist):
			continue
		default:
			return false, statErr
		}
	}

	return false, nil
}

// writePolicyFiles writes STYLE.md then quill.toml. Each file is created with an exclusive,
// non-clobbering open, so a file that appears between the occupancy check and the write (a
// check/write race or a concurrent process) fails instead of being truncated. If the quill.toml
// write fails, the already-written STYLE.md is removed so init never leaves a half-initialised
// repository.
func writePolicyFiles(
	stylePath string,
	profilePath string,
	styleGuide string,
	profile string,
) (err error) {
	if err = writeExclusive(stylePath, styleGuide); err != nil {
		return fmt.Errorf("write %s: %w", styleFileName, err)
	}
	if err = writeExclusive(profilePath, profile); err != nil {
		_ = os.Remove(stylePath)
		return fmt.Errorf("write %s: %w", profileFileName, err)
	}

	return nil
}

// writeExclusive atomically creates path containing contents with an exclusive open. It fails if
// path already exists and never truncates an existing file.
func writeExclusive(path string, contents string) (err error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, configFilePermission)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	if _, err = file.WriteString(contents); err != nil {
		return err
	}

	return nil
}

func parseInitOptionsWithResolver(
	_ repositoryRootResolver,
	arguments []string,
) (options initOptions, err error) {
	flagSet := newInitFlagSet(&options)
	if _, err = parseFlags(flagSet, initSummary, arguments); err != nil {
		return options, err
	}

	options.preset, err = parsePreset(options.preset)
	if err != nil {
		return options, err
	}

	options.repoRoot, err = resolveInitTarget(options.repoRoot)
	return options, err
}

// resolveInitTarget resolves the directory init writes into. Unlike repository-root discovery,
// init targets a directory that may not yet be a repository: an explicit path wins, otherwise the
// current working directory is used.
func resolveInitTarget(path string) (target string, err error) {
	if path != "" {
		return filepath.Abs(path)
	}

	return os.Getwd()
}

func newInitFlagSet(options *initOptions) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("init")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"target directory (current directory when omitted)",
	)
	flagSet.StringVar(&options.preset, "preset", defaultPreset, "preset: minimal")
	return flagSet
}

func initUsageText() (usage string) {
	var options initOptions
	return commandUsage("init", initSummary, newInitFlagSet(&options))
}

func initMachineMode(_ []string) (requested bool) {
	return false
}

/* ------------------------------------------- Presets ------------------------------------------ */

// minimalRequirementSection and minimalRequirementSlug name the single requirement the minimal
// preset documents and binds. They are kept as named parts rather than a combined literal so the
// generated requirement id lives only in the written STYLE.md and quill.toml, never in
// implementation source (requirement ids are repository-owned policy, bound in the profile).
const (
	minimalRequirementSection = "1.1"
	minimalRequirementSlug    = "enforcement-levels"
	requirementIDToken        = "__REQUIREMENT_ID__"
)

func minimalRequirementID() (id string) {
	return minimalRequirementSection + "." + minimalRequirementSlug
}

// initPreset returns the generated file contents for the named preset. Only "minimal" is
// supported; it is a concise, immediately-valid profile built on the built-in project Pack.
func initPreset(name string) (content initPresetContent, ok bool) {
	if name != defaultPreset {
		return initPresetContent{}, false
	}

	requirementID := minimalRequirementID()
	return initPresetContent{
		styleGuide: strings.ReplaceAll(minimalStyleGuide, requirementIDToken, requirementID),
		profile:    strings.ReplaceAll(minimalProfile, requirementIDToken, requirementID),
	}, true
}

// minimalStyleGuide is a small STYLE.md with one stable requirement bound by the minimal profile.
const minimalStyleGuide = `# Style

This STYLE.md documents the requirements Quill enforces. Each requirement carries a stable
identifier that the Quill profile binds to an executable rule. Replace this file with your
repository's real style policy as it grows.

## 1. Repository policy

### 1.1 Enforcement levels

<!-- style: id=__REQUIREMENT_ID__ -->
* Every bound style rule declares its enforcement level so failures classify consistently.
`

// minimalProfile enables the built-in project Pack, pins its required toolchain, and binds one
// representative, tool-free rule. It loads, compiles, and passes `quill check` without any
// installed tools. The project Pack requires every tool its rules reference to be pinned, so the
// seven canonical tools are pinned at their current versions; adjust them to suit your repository.
const minimalProfile = `# Quill profile generated by ` + "`quill init --preset minimal`" + `.
#
# A minimal, immediately-valid starter built on the built-in project Pack. Add scopes, Packs,
# tools, and rules as your repository grows. Inspect the active configuration with
# ` + "`quill list`" + ` and ` + "`quill explain rule:<id>`" + `.

schema_version = 1

[repository]
root_markers = ["STYLE.md", "quill.toml"]
default_scope = "all"
excluded_directories = [".git", "vendor"]
generated_marker = "DO NOT EDIT."

[repository.scope_roots]
all = ["."]

[style_guide]
path = "STYLE.md"

[packs]
enabled = ["project"]

# The project Pack requires a quality-command surface. Declare the Make targets your repository
# exposes; Quill validates them only when you bind the profile/quality-commands rule.
[packs.project.commands]
runner = "make"
path = "Makefile"

[[packs.project.commands.required_targets]]
name = "build"
recipe_line = "go build ./..."

# The project Pack references these tools, so each must be pinned. The bound rule below needs none
# of them, so ` + "`quill check`" + ` passes before they are installed.
[tools.go]
version = "1.24.5"

[tools.goimports]
version = "v0.42.0"

[tools.misspell]
version = "v0.3.4"

[tools."golangci-lint"]
version = "v2.6.2"

[tools.shfmt]
version = "v3.12.0"

[tools.shellcheck]
version = "0.10.0"

[tools.markdownlint]
version = "0.45.0"

[[rules]]
id = "profile/enforcement-levels"
enforcement = "required"
scope = "all"
requirement_ids = ["__REQUIREMENT_ID__"]
`

type initPresetContent struct {
	styleGuide string
	profile    string
}
