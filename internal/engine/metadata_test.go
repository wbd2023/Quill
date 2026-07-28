package engine

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
)

/* -------------------------------------------- Tests ------------------------------------------- */

// TestMetadataInspectsNoTools asserts the observable side effect of Metadata being metadata-only:
// it shares the prepared snapshot and the catalogue, never spawning a process or inspecting a tool.
func TestMetadataInspectsNoTools(t *testing.T) {
	t.Parallel()

	runner := &recordingRunner{}
	engine, err := New(testutil.RepositoryRoot(t), WithCommandRunner(runner))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if _, err = engine.Metadata(context.Background()); err != nil {
		t.Fatalf("Metadata: %v", err)
	}

	if runner.calls != 0 {
		t.Fatalf("Metadata invoked the command runner %d time(s), want 0", runner.calls)
	}
}

func TestMetadataRejectsCancelledContext(t *testing.T) {
	engine, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	operationContext, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err = engine.Metadata(operationContext); !errors.Is(err, context.Canceled) {
		t.Fatalf("Metadata error = %v, want context.Canceled", err)
	}
}

func TestMetadataReportsAvailableAndActive(t *testing.T) {
	t.Parallel()

	engine, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	snapshot, err := engine.Metadata(context.Background())
	if err != nil {
		t.Fatalf("Metadata: %v", err)
	}

	if len(snapshot.Packs) == 0 || len(snapshot.Rules) == 0 ||
		len(snapshot.Tools) == 0 || len(snapshot.Scopes) == 0 {
		t.Fatalf("metadata snapshot must populate packs, rules, tools, and scopes: %+v", snapshot)
	}

	var sawActivePack, sawActiveRule, sawDefaultScope bool
	for _, pack := range snapshot.Packs {
		if pack.External {
			t.Fatalf("shipped packs must be built-in, %q marked external", pack.ID)
		}
		if pack.Active {
			sawActivePack = true
		}
	}

	for _, rule := range snapshot.Rules {
		if rule.Check.Category == "" {
			t.Fatalf("rule %q missing check execution category", rule.ID)
		}
		if rule.Active {
			sawActiveRule = true
			if rule.Enforcement == "" || rule.Scope == "" || len(rule.RequirementIDs) == 0 {
				t.Fatalf("active rule %q missing binding: %+v", rule.ID, rule)
			}
		}
	}

	for _, scope := range snapshot.Scopes {
		if scope.Default {
			sawDefaultScope = true
		}
	}

	if !sawActivePack {
		t.Fatal("expected at least one active pack")
	}
	if !sawActiveRule {
		t.Fatal("expected at least one active rule")
	}
	if !sawDefaultScope {
		t.Fatal("expected a default scope")
	}
}

func TestMetadataIsDeterministic(t *testing.T) {
	t.Parallel()

	engine, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	first, err := engine.Metadata(context.Background())
	if err != nil {
		t.Fatalf("first Metadata: %v", err)
	}

	second, err := engine.Metadata(context.Background())
	if err != nil {
		t.Fatalf("second Metadata: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("Metadata must be deterministic across calls:\nfirst:  %+v\nsecond: %+v",
			first, second)
	}
}
