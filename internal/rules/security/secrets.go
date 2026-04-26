package security

import (
	"regexp"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
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
	code    string
	message string
	pattern *regexp.Regexp
}

/* --------------------------------------- Secret Scanning -------------------------------------- */

func CheckSecrets(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	patterns := committedSecretPatterns()

	files, err := filewalk.CollectAllFiles(repoRoot, repository, scope)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			for _, pattern := range patterns {
				if !pattern.pattern.MatchString(line.Text) {
					continue
				}

				result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
					Code:    pattern.code,
					File:    filewalk.RelativePath(repoRoot, path),
					Line:    line.Number,
					Message: "contains " + pattern.message,
				})
				break
			}

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

func committedSecretPatterns() (patterns []secretPattern) {
	return []secretPattern{
		{
			code:    "security/secrets/private-key",
			message: "possible private key material",
			pattern: regexp.MustCompile(pemMarkerPattern),
		},
		{
			code:    "security/secrets/aws-key",
			message: "possible AWS access key",
			pattern: regexp.MustCompile(awsIDPattern),
		},
		{
			code:    "security/secrets/github-pat",
			message: "possible GitHub personal access token",
			pattern: regexp.MustCompile(githubPersonalPattern),
		},
		{
			code:    "security/secrets/github-fine-grained-pat",
			message: "possible GitHub fine-grained token",
			pattern: regexp.MustCompile(githubFinePattern),
		},
		{
			code:    "security/secrets/slack-token",
			message: "possible Slack token",
			pattern: regexp.MustCompile(slackPattern),
		},
	}
}
