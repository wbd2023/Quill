package bindings

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/style"
)

/* ------------------------------------ Binding Completeness ------------------------------------ */

// TestEveryShippedExecutionIdentityHasExactlyOneRuntimeBinding is the binding-completeness contract
// required before co-locating Pack bindings. Every execution identity referenced by any shipped
// Rule's Check or Fix template must resolve to exactly one runtime binding produced by Build.
//
// "At most one" is enforced structurally: drivers.Bindings panics on duplicate registration, so a
// duplicate would fail Build itself. "At least one" is enforced here by resolving each identity
// through the public Bindings lookups. Together they prove exactly one binding per identity.
func TestEveryShippedExecutionIdentityHasExactlyOneRuntimeBinding(t *testing.T) {
	t.Parallel()

	registry, err := shipped.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	// Build itself enforces uniqueness: any duplicate registration panics inside
	// drivers.Bindings. If a duplicate slipped in, this would fail the test.
	built := Build()

	var missing []string
	for _, rule := range registry.Rules() {
		if rule.Check != nil {
			missing = append(missing, assertTemplateBound(rule.ID, "check", rule.Check, built)...)
		}
		if rule.Fix != nil {
			missing = append(missing, assertTemplateBound(rule.ID, "fix", rule.Fix, built)...)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf(
			"shipped rules reference %d unbound execution identit%s:\n  %s",
			len(missing),
			pluralY(missing),
			joinLines(missing),
		)
	}
}

// assertTemplateBound resolves one template's execution identity against built and returns a
// description (rule side, kind, key) for each identity that has no binding.
func assertTemplateBound(
	ruleID string,
	side string,
	template style.Template,
	built drivers.Bindings,
) (missing []string) {
	switch execution := template.(type) {
	case style.ToolchainExecution:
		// Toolchain execution has no per-identity binding; it always resolves to the shared
		// Toolchain driver. Nothing to verify.
		return nil

	case style.ProfileExecution:
		if _, found := built.LookupProfileCheck(execution.Check); !found {
			missing = append(missing, describe(ruleID, side, "profile check", execution.Check))
		}

	case style.FileCommandExecution:
		if _, found := built.LookupFileInterpreter(execution.ToolID); !found {
			missing = append(missing, describe(ruleID, side, "file interpreter", execution.ToolID))
		}

	case style.RepositoryScanExecution:
		if _, found := built.LookupRepositoryScanner(execution.Scanner); !found {
			missing = append(
				missing,
				describe(ruleID, side, "repository scanner", execution.Scanner),
			)
		}

	case style.TargetCommandTemplate:
		if _, found := built.LookupTargetCommand(execution.Action); !found {
			missing = append(missing, describe(ruleID, side, "target command", execution.Action))
		}

	case style.TargetCheckTemplate:
		if _, found := built.LookupTargetCheck(execution.Language); !found {
			missing = append(missing, describe(ruleID, side, "target check", execution.Language))
		}

	default:
		missing = append(
			missing,
			describe(ruleID, side, "unknown template", fmt.Sprintf("%T", template)),
		)
	}

	return missing
}

func describe(ruleID string, side string, kind string, key string) (description string) {
	return fmt.Sprintf("%s: %s %s %q", ruleID, side, kind, key)
}

func pluralY(items []string) (suffix string) {
	if len(items) == 1 {
		return "y"
	}
	return "ies"
}

func joinLines(items []string) (joined string) {
	var builder strings.Builder
	for _, item := range items {
		builder.WriteString("\n  - ")
		builder.WriteString(item)
	}
	return builder.String()
}
