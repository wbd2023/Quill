package vocabulary

import (
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
)

const goTypeSuffixMatchLength = 2
const goIdentifierSuffixMatchLength = 2
const bashAssignmentMatchLength = 4

func CheckVocabulary(
	repoRoot string,
	repository policy.RepositoryConfig,
	config Config,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	goFiles, err := filewalk.CollectFiles(repoRoot, repository, scope, ".go")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	shellFiles, err := filewalk.CollectFiles(repoRoot, repository, scope, ".sh")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range goFiles {
		err = checkGoVocabulary(&result, repoRoot, path, config)
		if err != nil {
			return style.ExecutionResult{}, err
		}
	}

	for _, path := range shellFiles {
		err = checkBashVocabulary(&result, repoRoot, path, config)
		if err != nil {
			return style.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, style.ViolationsFound()
}
