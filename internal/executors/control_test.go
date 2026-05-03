package executors

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckGlobalExclusionsPassesCurrentCollectorPolicy(t *testing.T) {
	if _, err := checkGlobalExclusions(profiles.Current(t).Repository); err != nil {
		t.Fatalf("checkGlobalExclusions: %v", err)
	}
}

func TestCheckQualitySurfaceAcceptsExpectedShape(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		filepath.Join("mk", "quality.mk"),
		`LINT_REQUIRED_ARGS = --mode required
LINT_FULL_ARGS = --mode all --strict-recommendations --verbose

lint: $(STYLE_BIN)
	@$(STYLE_CMD) check $(LINT_FULL_ARGS)

lint-required: $(STYLE_BIN)
	@$(STYLE_CMD) check $(LINT_REQUIRED_ARGS)

lint-fix: $(STYLE_BIN)
	@$(STYLE_CMD) fix --scope all

style-install: $(STYLE_BIN)
	@$(STYLE_CMD) install

style-doctor: $(STYLE_BIN)
	@$(STYLE_CMD) doctor

style-coverage: $(STYLE_BIN)
	@$(STYLE_CMD) coverage
	`,
	)

	if _, err := checkQualitySurface(repoRoot, profiles.Current(t).QualitySurface); err != nil {
		t.Fatalf("checkQualitySurface: %v", err)
	}
}

func TestCheckQualitySurfaceRejectsMissingRequiredRecipe(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		filepath.Join("mk", "quality.mk"),
		`LINT_REQUIRED_ARGS = --mode required
LINT_FULL_ARGS = --mode all --strict-recommendations --verbose
lint: $(STYLE_BIN)
	@$(STYLE_CMD) check $(LINT_FULL_ARGS)
	`,
	)

	if _, err := checkQualitySurface(repoRoot, profiles.Current(t).QualitySurface); err == nil {
		t.Fatal("expected missing lint-required recipe to fail")
	}
}
