package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
)

type trackingPackProvider struct {
	defaultPackProvider
	runtimeCalls int
}

func (provider *trackingPackProvider) Runtime(
	operationContext context.Context,
	enabledPacks []string,
) (runtime PackRuntime, loadError error) {
	provider.runtimeCalls++
	return provider.defaultPackProvider.Runtime(operationContext, enabledPacks)
}

func TestCoverageDoesNotConstructPackRuntime(t *testing.T) {
	provider := &trackingPackProvider{}
	engine, err := New(
		testutil.RepositoryRoot(t),
		WithPackProvider(provider),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if _, err = engine.Coverage(context.Background()); err != nil {
		t.Fatalf("Coverage: %v", err)
	}

	if provider.runtimeCalls != 0 {
		t.Fatalf("Pack runtime calls = %d, want 0", provider.runtimeCalls)
	}
}

func TestCoverageRejectsCancelledContext(t *testing.T) {
	engine, err := New(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	operationContext, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err = engine.Coverage(operationContext); !errors.Is(err, context.Canceled) {
		t.Fatalf("Coverage error = %v, want context.Canceled", err)
	}
}
