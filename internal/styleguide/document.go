package styleguide

import (
	"os"
	"path/filepath"
)

const (
	VerificationAutomated      VerificationMode = "automated"
	VerificationReviewOnly     VerificationMode = "review_only"
	VerificationManualDeferred VerificationMode = "manual_deferred"
)

type VerificationMode string

type Config struct {
	Path                string
	RequirementIDFormat string
}

type Document struct {
	Headings     []Heading
	Requirements []Requirement
}

type Heading struct {
	Section string
	Title   string
}

type Requirement struct {
	ID      string
	Section string
	Text    string
	Mode    VerificationMode
	Reason  string
}

type RequirementMetadata struct {
	ID     string
	Mode   VerificationMode
	Reason string
}

func Load(repoRoot string, config Config) (document Document, err error) {
	stylePath := filepath.Join(repoRoot, config.Path)
	contents, err := os.ReadFile(stylePath)
	if err != nil {
		return Document{}, err
	}

	return Parse(contents, config)
}

func Parse(contents []byte, config Config) (document Document, err error) {
	return compileDocument(contents, config)
}
