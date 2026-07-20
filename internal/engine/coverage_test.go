package engine

import (
	"context"
	"testing"

	"github.com/wbd2023/Quill/internal/testutil"
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
