package vocabulary

import (
	"github.com/wbd2023/Quill/internal/checks/vocabularypolicy"
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

const goTypeSuffixMatchLength = 2
const goIdentifierSuffixMatchLength = 2
const bashAssignmentMatchLength = 4

// CheckVocabulary check vocabulary.
func CheckVocabulary(
	repoRoot string,
	repository policy.RepositoryConfig,
	config vocabularypolicy.Config,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	roots := repository.ResolveScopeRoots(repoRoot, scope)
	walkConfig := filewalk.WalkConfig{
		ExcludedDirectories: repository.ExcludedDirectories,
		GeneratedMarker:     repository.GeneratedMarker,
	}

	goFiles, err := filewalk.CollectFiles(roots, walkConfig, ".go")
	if err != nil {
		return style.ExecutionResult{}, err
	}

	shellFiles, err := filewalk.CollectFiles(roots, walkConfig, ".sh")
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

	return result, nil
}
