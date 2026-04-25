package repostyle

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
)

/* ---------------------------------------- Naming Rules ---------------------------------------- */

const goTypeSuffixMatchLength = 2
const goIdentifierSuffixMatchLength = 2
const shellAssignmentMatchLength = 4

func CheckNaming(
	repoRoot string,
	repository profile.RepositoryConfig,
	naming profile.NamingConfig,
	scope contract.Scope,
) (output string, err error) {
	goTypePattern := compileGoTypeSuffixPattern(naming.GoTypeSuffixForbidden)
	goIdentifierPattern := compileGoIdentifierSuffixPattern(naming.GoIdentifierSuffixForbidden)
	shellAssignmentPattern := compileShellAssignmentPattern(naming.ShellForbiddenAssignments)

	goFiles, err := CollectFiles(repoRoot, repository, scope, ".go")
	if err != nil {
		return "", err
	}

	shellFiles, err := CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	found := false

	for _, path := range goFiles {
		file, openErr := os.Open(path)
		if openErr != nil {
			return "", openErr
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if suffix := matchedGoTypeSuffix(goTypePattern, line); suffix != "" {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d use %s not %s in type names\n",
					RelativePath(repoRoot, path),
					lineNumber,
					naming.GoTypeSuffixPreferred,
					suffix,
				))
			}

			if naming.GoIdentifierSuffixPreferred != "" &&
				strings.Contains(line, naming.GoIdentifierSuffixPreferred) {
				continue
			}

			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}

			if suffix := matchedGoIdentifierSuffix(goIdentifierPattern, line); suffix != "" {
				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d use x%s not x%s\n",
					RelativePath(repoRoot, path),
					lineNumber,
					naming.GoIdentifierSuffixPreferred,
					suffix,
				))
			}
		}

		if scanErr := scanner.Err(); scanErr != nil {
			return "", closeFile(file, scanErr)
		}

		if closeErr := closeFile(file, nil); closeErr != nil {
			return "", closeErr
		}
	}

	for _, path := range shellFiles {
		file, openErr := os.Open(path)
		if openErr != nil {
			return "", openErr
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			name := matchedShellAssignment(shellAssignmentPattern, line)
			if name == "" {
				continue
			}

			found = true
			builder.WriteString(fmt.Sprintf(
				"%s:%d use descriptive constant names in Bash "+
					"(prefer %s over %s)\n",
				RelativePath(repoRoot, path),
				lineNumber,
				naming.ShellPreferredAssignment,
				name,
			))
		}

		if scanErr := scanner.Err(); scanErr != nil {
			return "", closeFile(file, scanErr)
		}

		if closeErr := closeFile(file, nil); closeErr != nil {
			return "", closeErr
		}
	}

	if !found {
		return "", nil
	}

	return builder.String(), errViolationsFound
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
