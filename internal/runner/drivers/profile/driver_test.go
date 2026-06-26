package profile

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/checks/projectpolicy"
	projectpack "ciphera/tools/internal/pack/shipped/project"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

func TestCheckExcludedDirectoriesPassesCurrentCollectorPolicy(t *testing.T) {
	if _, err := checkExcludedDirectories(profiles.Current(t).Repository); err != nil {
		t.Fatalf("checkExcludedDirectories: %v", err)
	}
}

func TestCheckCommandsAcceptsExpectedShape(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
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

	if _, err := checkCommands(repoRoot, currentCommands(t)); err != nil {
		t.Fatalf("checkCommands: %v", err)
	}
}

func TestCheckCommandsRejectsMissingRequiredRecipe(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
		t,
		repoRoot,
		filepath.Join("mk", "quality.mk"),
		`LINT_REQUIRED_ARGS = --mode required
LINT_FULL_ARGS = --mode all --strict-recommendations --verbose
lint: $(STYLE_BIN)
	@$(STYLE_CMD) check $(LINT_FULL_ARGS)
	`,
	)

	if _, err := checkCommands(repoRoot, currentCommands(t)); err == nil {
		t.Fatal("expected missing lint-required recipe to fail")
	}
}

func currentCommands(t *testing.T) (commands projectpolicy.CommandsConfig) {
	t.Helper()

	pack, found := profiles.Current(t).PackConfigs.Lookup(projectpack.PackID)
	if !found {
		t.Fatal("missing project pack config")
	}

	config, err := projectpolicy.DecodeConfig(pack)
	if err != nil {
		t.Fatalf("Decode project config: %v", err)
	}

	return config.Commands
}
