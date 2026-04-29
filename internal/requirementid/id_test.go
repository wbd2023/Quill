package requirementid_test

import (
	"testing"

	"ciphera/tools/internal/requirementid"
)

/* ------------------------------------------- Parsing ------------------------------------------ */

func TestParse(t *testing.T) {
	id, err := requirementid.Parse(
		"3.8.constructor-category-order",
		requirementid.SectionSlug,
	)
	if err != nil {
		t.Fatalf("parse requirement id: %v", err)
	}

	if id.String() != "3.8.constructor-category-order" {
		t.Fatalf("unexpected requirement id %q", id.String())
	}
	if id.Section() != "3.8" {
		t.Fatalf("unexpected requirement section %q", id.Section())
	}
	if id.Slug() != "constructor-category-order" {
		t.Fatalf("unexpected requirement slug %q", id.Slug())
	}
}

func TestParseRejectsInvalidID(t *testing.T) {
	cases := []struct {
		name  string
		value string
	}{
		{name: "missing section", value: "3"},
		{name: "missing slug", value: "3.8"},
		{name: "empty slug", value: "3.8."},
		{name: "leading hyphen", value: "3.8.-constructor-order"},
		{name: "repeated hyphen", value: "3.8.constructor--order"},
		{name: "trailing hyphen", value: "3.8.constructor-order-"},
		{name: "uppercase slug", value: "3.8.Constructor-order"},
		{name: "invalid section", value: "3.x.constructor-order"},
		{name: "leading zero major", value: "03.8.constructor-order"},
		{name: "leading zero minor", value: "3.08.constructor-order"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := requirementid.Parse(test.value, requirementid.SectionSlug)
			if err == nil {
				t.Fatalf("expected requirement id %q to be rejected", test.value)
			}
		})
	}
}

func TestParseRejectsUnsupportedScheme(t *testing.T) {
	_, err := requirementid.Parse("3.8.constructor-category-order", "unknown")
	if err == nil {
		t.Fatalf("expected unsupported requirement id scheme to be rejected")
	}
}

/* ------------------------------------------ Sections ------------------------------------------ */

func TestValidSection(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected bool
	}{
		{name: "valid section", value: "3.8", expected: true},
		{name: "zero major", value: "0.1", expected: true},
		{name: "two digit minor", value: "3.10", expected: true},
		{name: "empty section", value: "", expected: false},
		{name: "missing minor", value: "3", expected: false},
		{name: "empty major", value: ".8", expected: false},
		{name: "empty minor", value: "3.", expected: false},
		{name: "invalid major", value: "x.8", expected: false},
		{name: "invalid minor", value: "3.x", expected: false},
		{name: "leading zero major", value: "03.8", expected: false},
		{name: "leading zero minor", value: "3.08", expected: false},
		{name: "too many parts", value: "3.8.1", expected: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			valid := requirementid.ValidSection(test.value)
			if valid != test.expected {
				t.Fatalf("unexpected section validity %t", valid)
			}
		})
	}
}
