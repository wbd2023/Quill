package external

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

/* --------------------------------------- Source Loading --------------------------------------- */

// ManifestName is the manifest filename read from a Pack source directory.
const ManifestName = "pack.toml"

// defaultGroup is assigned to an external rule that declares no group.
const defaultGroup style.RuleGroup = "external"

// LoadSources loads every declared local Pack source beneath repoRoot into normal Pack definitions.
// Each source directory must contain a valid pack.toml and a resolvable runtime executable;
// both are verified before the definition is returned, so no external Pack process may launch
// until its manifest and executable pass validation.
//
// repoRoot must be canonical. Source paths are joined beneath it and the resulting Pack directory
// is verified to remain inside the repository after symlink resolution.
func LoadSources(
	repoRoot string,
	sources []profile.PackSource,
) (definitions []pack.Definition, err error) {
	canonicalRoot, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve repository root for pack sources: %w", err)
	}

	definitions = make([]pack.Definition, 0, len(sources))
	for _, source := range sources {
		definition, err := loadSource(canonicalRoot, source)
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func loadSource(
	canonicalRoot string,
	source profile.PackSource,
) (definition pack.Definition, err error) {
	packDirectory := filepath.Join(canonicalRoot, filepath.Clean(source.Path))
	resolved, err := filepath.EvalSymlinks(packDirectory)
	if err != nil {
		return pack.Definition{}, fmt.Errorf("resolve pack source %q: %w", source.Path, err)
	}
	if !isWithin(canonicalRoot, resolved) {
		return pack.Definition{}, fmt.Errorf(
			"pack source %q escapes the repository root",
			source.Path,
		)
	}

	manifestPath, err := resolveManifest(resolved)
	if err != nil {
		return pack.Definition{}, err
	}

	contents, err := os.ReadFile(manifestPath)
	if err != nil {
		return pack.Definition{}, fmt.Errorf("read pack manifest %q: %w", manifestPath, err)
	}

	manifest, err := DecodeManifest(string(contents))
	if err != nil {
		return pack.Definition{}, err
	}

	if _, err = ResolveExecutable(resolved, manifest.Runtime.Command); err != nil {
		return pack.Definition{}, err
	}

	return toDefinition(manifest, resolved), nil
}

func toDefinition(manifest Manifest, packDirectory string) (definition pack.Definition) {
	rules := make([]style.RuleDefinition, 0, len(manifest.Rules))
	for _, rule := range manifest.Rules {
		group := rule.Group
		if group == "" {
			group = string(defaultGroup)
		}

		name := rule.Name
		if name == "" {
			name = rule.ID
		}

		rules = append(rules, style.RuleDefinition{
			ID:    rule.ID,
			Name:  name,
			Group: style.RuleGroup(group),
			Check: style.ExternalCheck{
				CheckID:       rule.Check,
				FileSet:       rule.FileSet,
				PackDirectory: packDirectory,
				Command:       manifest.Runtime.Command,
				Timeout:       manifest.Runtime.Timeout,
			},
		})
	}

	return pack.Definition{
		ID:    manifest.Pack.ID,
		Name:  manifest.Pack.Name,
		Rules: rules,
		Policy: pack.Policy{
			Validate: acceptPackPolicy,
		},
	}
}

// acceptPackPolicy accepts an external Pack Policy unchanged. It is forwarded to the subprocess
// through Request.Policy, so the external Pack validates its own policy rather than Quill
// interpreting Pack-specific values.
func acceptPackPolicy(_ profile.PackPolicy) (err error) {
	return nil
}

// resolveManifest locates pack.toml beneath a canonical Pack directory and resolves it through
// symlinks before returning the physical path. The resolved manifest must stay inside the Pack
// directory and be a regular file: a symlinked manifest pointing outside the Pack, or a non-regular
// file, is rejected so untrusted repository input cannot redirect the manifest read.
func resolveManifest(packDirectory string) (manifestPath string, err error) {
	joined := filepath.Join(packDirectory, ManifestName)
	resolved, err := filepath.EvalSymlinks(joined)
	if err != nil {
		return "", fmt.Errorf("resolve pack manifest %q: %w", ManifestName, err)
	}
	if !isWithin(packDirectory, resolved) {
		return "", fmt.Errorf("pack manifest %q resolves outside the Pack directory", ManifestName)
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("resolve pack manifest %q: %w", ManifestName, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("pack manifest %q must be a regular file", ManifestName)
	}

	return resolved, nil
}

/* ----------------------------------------- Containment ---------------------------------------- */

// isWithin reports whether target is root or a descendant of root. Both arguments must be absolute
// and cleaned. It is the containment check shared by source loading (beneath the repository) and
// executable resolution (beneath the Pack directory).
func isWithin(root string, target string) (within bool) {
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	cleaned := filepath.Clean(relative)
	return cleaned != ".." && !strings.HasPrefix(cleaned, ".."+string(filepath.Separator))
}
