package naming

import (
	"fmt"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

/* ---------------------------------------- Naming Rules ---------------------------------------- */

const goTypeSuffixMatchLength = 2
const goIdentifierSuffixMatchLength = 2
const shellAssignmentMatchLength = 4

func CheckNaming(
	repoRoot string,
	repository policy.RepositoryConfig,
	naming policy.NamingConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	goTypePattern := compileGoTypeSuffixPattern(naming.GoTypeSuffixForbidden)
	goIdentifierPattern := compileGoIdentifierSuffixPattern(naming.GoIdentifierSuffixForbidden)
	shellAssignmentPattern := compileShellAssignmentPattern(naming.ShellForbiddenAssignments)

	goFiles, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	shellFiles, err := filewalk.CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range goFiles {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if suffix := matchedGoTypeSuffix(goTypePattern, line.Text); suffix != "" {
				result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
					Code: "naming/vocabulary/go-type-suffix",
					File: filewalk.RelativePath(repoRoot, path),
					Line: line.Number,
					Message: fmt.Sprintf(
						"use %s not %s in type names",
						naming.GoTypeSuffixPreferred,
						suffix,
					),
				})
			}

			if naming.GoIdentifierSuffixPreferred != "" &&
				strings.Contains(line.Text, naming.GoIdentifierSuffixPreferred) {
				return nil
			}

			if strings.HasPrefix(strings.TrimSpace(line.Text), "//") {
				return nil
			}

			if suffix := matchedGoIdentifierSuffix(goIdentifierPattern, line.Text); suffix != "" {
				result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
					Code: "naming/vocabulary/go-identifier-suffix",
					File: filewalk.RelativePath(repoRoot, path),
					Line: line.Number,
					Message: fmt.Sprintf(
						"use x%s not x%s",
						naming.GoIdentifierSuffixPreferred,
						suffix,
					),
				})
			}

			return nil
		})
		if err != nil {
			return contract.ExecutionResult{}, err
		}
	}

	for _, path := range shellFiles {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			name := matchedShellAssignment(shellAssignmentPattern, line.Text)
			if name == "" {
				return nil
			}

			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code: "naming/vocabulary/shell-assignment",
				File: filewalk.RelativePath(repoRoot, path),
				Line: line.Number,
				Message: fmt.Sprintf(
					"use descriptive constant names in Bash (prefer %s over %s)",
					naming.ShellPreferredAssignment,
					name,
				),
			})
			return nil
		})
		if err != nil {
			return contract.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}

func compileGoTypeSuffixPattern(suffixes []string) (pattern *regexp.Regexp) {
	if len(suffixes) == 0 {
		return nil
	}

	return regexp.MustCompile(
		fmt.Sprintf(`type\s+\w*(%s)\s+`, strings.Join(suffixes, "|")),
	)
}

func compileGoIdentifierSuffixPattern(suffixes []string) (pattern *regexp.Regexp) {
	if len(suffixes) == 0 {
		return nil
	}

	return regexp.MustCompile(
		fmt.Sprintf(`\b\w+(%s)\b`, strings.Join(suffixes, "|")),
	)
}

func compileShellAssignmentPattern(names []string) (pattern *regexp.Regexp) {
	if len(names) == 0 {
		return nil
	}

	return regexp.MustCompile(
		fmt.Sprintf(`(^|[[:space:]])(local[[:space:]]+)?(%s)=`, strings.Join(names, "|")),
	)
}

func matchedGoTypeSuffix(pattern *regexp.Regexp, line string) (suffix string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < goTypeSuffixMatchLength {
		return ""
	}

	return matches[1]
}

func matchedGoIdentifierSuffix(pattern *regexp.Regexp, line string) (suffix string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < goIdentifierSuffixMatchLength {
		return ""
	}

	return matches[1]
}

func matchedShellAssignment(pattern *regexp.Regexp, line string) (name string) {
	if pattern == nil {
		return ""
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < shellAssignmentMatchLength {
		return ""
	}

	return matches[3]
}
