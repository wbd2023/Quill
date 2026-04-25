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

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	pemMarkerPattern = `-----BEG` +
		`IN (?:[A-Z0-9]+(?: [A-Z0-9]+)*) PRI` +
		`VATE` +
		` KEY-----`
	awsIDPattern          = `\bAKI` + `A[0-9A-Z]{16}\b`
	githubPersonalPattern = `\bgh` + `p_[A-Za-z0-9]{36}\b`
	githubFinePattern     = `\bgithub_` + `pat_[A-Za-z0-9_]{20,}\b`
	slackPattern          = `\bxo` + `x` + `[baprs]-[A-Za-z0-9-]{10,}\b`
)

/* -------------------------------------------- Types ------------------------------------------- */

type secretPattern struct {
	message string
	pattern *regexp.Regexp
}

/* --------------------------------------- Secret Scanning -------------------------------------- */

func CheckSecrets(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	patterns := committedSecretPatterns()

	files, err := CollectAllFiles(repoRoot, repository, scope)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	found := false

	for _, path := range files {
		file, openErr := os.Open(path)
		if openErr != nil {
			return "", openErr
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()

			for _, pattern := range patterns {
				if !pattern.pattern.MatchString(line) {
					continue
				}

				found = true
				builder.WriteString(fmt.Sprintf(
					"%s:%d contains %s\n",
					RelativePath(repoRoot, path),
					lineNumber,
					pattern.message,
				))
				break
			}
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

func committedSecretPatterns() (patterns []secretPattern) {
	return []secretPattern{
		{
			message: "possible private key material",
			pattern: regexp.MustCompile(pemMarkerPattern),
		},
		{
			message: "possible AWS access key",
			pattern: regexp.MustCompile(awsIDPattern),
		},
		{
			message: "possible GitHub personal access token",
			pattern: regexp.MustCompile(githubPersonalPattern),
		},
		{
			message: "possible GitHub fine-grained token",
			pattern: regexp.MustCompile(githubFinePattern),
		},
		{
			message: "possible Slack token",
			pattern: regexp.MustCompile(slackPattern),
		},
	}
}
