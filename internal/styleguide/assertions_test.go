package styleguide

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func requireDocument(t *testing.T, document Document, expected Document) {
	t.Helper()

	if diff := cmp.Diff(expected, document); diff != "" {
		t.Fatalf("unexpected document (-expected +actual):\n%s", diff)
	}
}

func requireHeading(t *testing.T, heading Heading, expected Heading) {
	t.Helper()

	if diff := cmp.Diff(expected, heading); diff != "" {
		t.Fatalf("unexpected heading (-expected +actual):\n%s", diff)
	}
}

func requireHeadings(t *testing.T, headings []Heading, expected []Heading) {
	t.Helper()

	if diff := cmp.Diff(expected, headings); diff != "" {
		t.Fatalf("unexpected headings (-expected +actual):\n%s", diff)
	}
}

func requireRequirements(t *testing.T, requirements []Requirement, expected []Requirement) {
	t.Helper()

	if diff := cmp.Diff(expected, requirements); diff != "" {
		t.Fatalf("unexpected requirements (-expected +actual):\n%s", diff)
	}
}

func requireMetadataFields(t *testing.T, fields metadataFields, expected metadataFields) {
	t.Helper()

	if fields != expected {
		t.Fatalf("unexpected metadata fields\nexpected: %#v\nactual:   %#v", expected, fields)
	}
}

func requireErrorContains(t *testing.T, err error, expected string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected error containing %q, got %v", expected, err)
	}
}
