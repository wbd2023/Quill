package styleguide

import (
	"os"
	"path/filepath"

	"ciphera/tools/internal/profile"
)

/* -------------------------------------------- Types ------------------------------------------- */

type documentHeading struct {
	Section string
	Title   string
}

type documentRequirement struct {
	ID      string
	Section string
	Text    string
	Mode    VerificationMode
	Reason  string
}

type documentModel struct {
	Headings     []documentHeading
	Requirements []documentRequirement
}

/* -------------------------------------- Document Loading -------------------------------------- */

func readHeadings(repoRoot string) (headings []documentHeading, err error) {
	document, err := readDocument(repoRoot)
	if err != nil {
		return nil, err
	}

	return document.Headings, nil
}

func readRequirements(repoRoot string) (requirements []documentRequirement, err error) {
	document, err := readDocument(repoRoot)
	if err != nil {
		return nil, err
	}

	return document.Requirements, nil
}

func readDocument(repoRoot string) (document documentModel, err error) {
	policy, source, err := readStyleGuide(repoRoot)
	if err != nil {
		return documentModel{}, err
	}

	return compileDocument(source, policy.StyleGuide)
}

func readStyleGuide(
	repoRoot string,
) (policy profile.Profile, contents []byte, err error) {
	policy, err = profile.Load(repoRoot)
	if err != nil {
		return profile.Profile{}, nil, err
	}

	stylePath := filepath.Join(repoRoot, policy.StyleGuide.Path)
	contents, err = os.ReadFile(stylePath)
	if err != nil {
		return profile.Profile{}, nil, err
	}

	return policy, contents, nil
}
