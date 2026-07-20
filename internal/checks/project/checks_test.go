package project

import (
	"testing"

	"github.com/wbd2023/Quill/internal/checks/projectpolicy"
	projectpack "github.com/wbd2023/Quill/internal/pack/shipped/project"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func TestCheckExcludedDirectoriesPassesCurrentCollectorPolicy(t *testing.T) {
	if _, err := CheckExcludedDirectories(profiles.Current(t).Repository); err != nil {
		t.Fatalf("CheckExcludedDirectories: %v", err)
	}
}

func TestCheckCommandsAcceptsExpectedShape(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
		t,
		repoRoot,
		"Makefile",
		`LINT_REQUIRED_ARGS = --mode required
LINT_FULL_ARGS = --mode all --strict-recommendations --verbose

lint:
	@$(QUILL_CMD) check $(LINT_FULL_ARGS)

lint-required:
	@$(QUILL_CMD) check $(LINT_REQUIRED_ARGS)

lint-fix:
	@$(QUILL_CMD) fix --scope all

style-install:
	@$(QUILL_CMD) install

style-doctor:
	@$(QUILL_CMD) doctor

style-coverage:
	@$(QUILL_CMD) coverage
	`,
	)

	if _, err := CheckCommands(repoRoot, currentCommands(t)); err != nil {
		t.Fatalf("CheckCommands: %v", err)
	}
}

func TestCheckCommandsRejectsMissingRequiredRecipe(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
		t,
		repoRoot,
		"Makefile",
		`LINT_REQUIRED_ARGS = --mode required
LINT_FULL_ARGS = --mode all --strict-recommendations --verbose
lint:
	@$(QUILL_CMD) check $(LINT_FULL_ARGS)
	`,
	)

	if output, _ := CheckCommands(repoRoot, currentCommands(t)); output == "" {
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
