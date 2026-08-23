package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/engine"
)

/* ------------------------------------------ List JSON ----------------------------------------- */

func TestListJSONCarriesOnlySelectedSection(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		selector string
		result   ListResult
		wantKey  string
	}{
		{
			name:     "rules selector carries rules only",
			selector: ListRules,
			result: ListResult{
				Selector: ListRules,
				Packs:    []ListPack{{ID: "p", Name: "P", Active: true}},
				Rules: []ListRule{{
					ID: "r", Pack: "p", Name: "R", Active: true, Enforcement: "required",
				}},
			},
			wantKey: "rules",
		},
		{
			name:     "packs selector carries packs only",
			selector: ListPacks,
			result: ListResult{
				Selector: ListPacks,
				Packs:    []ListPack{{ID: "p", Name: "P"}},
				Rules:    []ListRule{{ID: "r", Pack: "p"}},
			},
			wantKey: "packs",
		},
	}

	for _, command := range cases {
		t.Run(command.name, func(t *testing.T) {
			t.Parallel()

			var buffer bytes.Buffer
			if err := WriteList(
				&buffer, testEnvelopeMetadata("list"), FormatJSON, command.result,
			); err != nil {
				t.Fatalf("WriteList: %v", err)
			}

			var envelope struct {
				Result map[string]json.RawMessage `json:"result"`
			}
			if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
				t.Fatalf("decode envelope: %v\n%s", err, buffer.String())
			}

			if _, found := envelope.Result[command.wantKey]; !found {
				t.Fatalf("result missing selected key %q: %s", command.wantKey, envelope.Result)
			}

			for _, other := range []string{ListPacks, ListRules, ListTools, ListScopes} {
				if other == command.wantKey {
					continue
				}
				if _, found := envelope.Result[other]; found {
					t.Fatalf("result must not carry unselected key %q: %s", other, envelope.Result)
				}
			}
		})
	}
}

func TestListRuleJSONOmitsEmptyEnforcementForInactiveRules(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	result := ListResult{
		Selector: ListRules,
		Rules: []ListRule{
			{
				ID: "active", Pack: "p", Name: "A", Active: true,
				Enforcement: "required", Scope: "all",
			},
			{ID: "inactive", Pack: "p", Name: "I", Active: false},
		},
	}

	if err := WriteList(&buffer, testEnvelopeMetadata("list"), FormatJSON, result); err != nil {
		t.Fatalf("WriteList: %v", err)
	}

	var envelope struct {
		Result struct {
			Rules []map[string]any `json:"rules"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}

	rules := envelope.Result.Rules
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}

	for _, rule := range rules {
		active, _ := rule["active"].(bool)
		if active {
			if _, hasEnforcement := rule["enforcement"]; !hasEnforcement {
				t.Fatalf("active rule must carry enforcement: %+v", rule)
			}
		} else {
			if _, hasEnforcement := rule["enforcement"]; hasEnforcement {
				t.Fatalf("inactive rule must omit enforcement: %+v", rule)
			}
		}
	}
}

func TestListPacksCarriesPackProvenance(t *testing.T) {
	t.Parallel()

	result := NewListResult(
		engine.MetadataSnapshot{
			Packs: []engine.PackMetadata{{
				ID:         "company",
				Name:       "Company",
				Provenance: engine.PackProvenanceExternal,
			}},
		},
		ListPacks,
	)

	var buffer bytes.Buffer
	if err := WriteList(&buffer, testEnvelopeMetadata("list"), FormatJSON, result); err != nil {
		t.Fatalf("WriteList: %v", err)
	}

	var envelope struct {
		Result struct {
			Packs []ListPack `json:"packs"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode envelope: %v\n%s", err, buffer.String())
	}
	if got, want := envelope.Result.Packs[0].Provenance, "external"; got != want {
		t.Fatalf("pack provenance = %q, want %q", got, want)
	}
}

func TestListRulesCarryOwningPackProvenance(t *testing.T) {
	t.Parallel()

	result := NewListResult(
		engine.MetadataSnapshot{
			Packs: []engine.PackMetadata{
				{ID: "shipped", Provenance: engine.PackProvenanceShipped},
				{ID: "external", Provenance: engine.PackProvenanceExternal},
			},
			Rules: []engine.RuleMetadata{
				{ID: "shipped/rule", PackID: "shipped"},
				{ID: "external/rule", PackID: "external"},
			},
		},
		ListRules,
	)

	if got, want := result.Rules[0].Provenance, "shipped"; got != want {
		t.Fatalf("Shipped Rule provenance = %q, want %q", got, want)
	}
	if got, want := result.Rules[1].Provenance, "external"; got != want {
		t.Fatalf("External Rule provenance = %q, want %q", got, want)
	}
}

func TestListToolJSONDoesNotImplyUnsupportedProvenance(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	result := ListResult{
		Selector: ListTools,
		Tools:    []ListTool{{ID: "tool", Name: "Tool", Packs: []string{"pack"}}},
	}
	if err := WriteList(&buffer, testEnvelopeMetadata("list"), FormatJSON, result); err != nil {
		t.Fatalf("WriteList: %v", err)
	}

	var envelope struct {
		Result struct {
			Tools []map[string]any `json:"tools"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode envelope: %v\n%s", err, buffer.String())
	}
	if _, found := envelope.Result.Tools[0]["external"]; found {
		t.Fatalf("Tool must not carry unsupported provenance: %+v", envelope.Result.Tools[0])
	}
}

/* ---------------------------------------- Explain JSON ---------------------------------------- */

func TestExplainJSONOmitsFixWhenAbsent(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	result := ExplainResult{
		Rule: ExplainRule{
			ID:   "profile/example",
			Name: "Example",
			Pack: ExplainPack{ID: "project"},
			Binding: ExplainBinding{
				Enforcement:  "required",
				Scope:        "all",
				Requirements: []string{"1.1.example"},
			},
			Check: ExplainExecution{Category: "profile"},
		},
	}

	if err := WriteExplain(
		&buffer, testEnvelopeMetadata("explain"), FormatJSON, result,
	); err != nil {
		t.Fatalf("WriteExplain: %v", err)
	}

	var envelope struct {
		Result struct {
			Rule map[string]any `json:"rule"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode: %v\n%s", err, buffer.String())
	}

	if _, hasFix := envelope.Result.Rule["fix"]; hasFix {
		t.Fatalf("explain must omit fix when absent: %+v", envelope.Result.Rule)
	}
}

func TestExplainJSONIncludesFixWhenPresent(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	fix := ExplainExecution{Category: "file_command"}
	result := ExplainResult{
		Rule: ExplainRule{
			ID:   "text/example",
			Name: "Example",
			Pack: ExplainPack{ID: "text"},
			Binding: ExplainBinding{
				Enforcement:  "required",
				Scope:        "all",
				Requirements: []string{"1.1.example"},
			},
			Check: ExplainExecution{Category: "file_command"},
			Fix:   &fix,
		},
	}

	if err := WriteExplain(
		&buffer, testEnvelopeMetadata("explain"), FormatJSON, result,
	); err != nil {
		t.Fatalf("WriteExplain: %v", err)
	}

	var envelope struct {
		Result struct {
			Rule struct {
				Fix *ExplainExecution `json:"fix"`
			} `json:"rule"`
		} `json:"result"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode: %v\n%s", err, buffer.String())
	}

	if envelope.Result.Rule.Fix == nil {
		t.Fatal("explain must include fix when present")
	}
}

/* ------------------------------------- Explain Pack Policy ------------------------------------ */

func TestExplainTextRendersPackPolicyDeterministically(t *testing.T) {
	t.Parallel()

	// Insert keys out of sorted order to prove the renderer emits them sorted.
	config := map[string]any{
		"zebra": "last",
		"alpha": map[string]any{
			"inner2": float64(2),
			"inner1": float64(1),
		},
	}

	result := ExplainResult{
		Rule: ExplainRule{
			ID:   "profile/example",
			Name: "Example",
			Pack: ExplainPack{
				ID:     "project",
				Policy: config,
			},
			Check: ExplainExecution{Category: "profile"},
		},
	}

	var buffer bytes.Buffer
	if err := WriteExplain(
		&buffer, testEnvelopeMetadata("explain"), FormatText, result,
	); err != nil {
		t.Fatalf("WriteExplain: %v", err)
	}

	output := buffer.String()
	if !strings.Contains(output, "Pack policy") {
		t.Fatalf("text explain must render a Pack Policy section:\n%s", output)
	}

	alphaIdx := strings.Index(output, `"alpha"`)
	zebraIdx := strings.Index(output, `"zebra"`)
	inner1Idx := strings.Index(output, `"inner1"`)
	inner2Idx := strings.Index(output, `"inner2"`)

	if alphaIdx < 0 || zebraIdx < 0 || alphaIdx > zebraIdx {
		t.Fatalf("policy keys must render sorted (alpha before zebra):\n%s", output)
	}
	if inner1Idx < 0 || inner2Idx < 0 || inner1Idx > inner2Idx {
		t.Fatalf("nested policy keys must render sorted (inner1 before inner2):\n%s", output)
	}
}

func TestExplainTextOmitsPackPolicyWhenEmpty(t *testing.T) {
	t.Parallel()

	result := ExplainResult{
		Rule: ExplainRule{
			ID:    "profile/example",
			Name:  "Example",
			Pack:  ExplainPack{ID: "project"},
			Check: ExplainExecution{Category: "profile"},
		},
	}

	var buffer bytes.Buffer
	if err := WriteExplain(
		&buffer, testEnvelopeMetadata("explain"), FormatText, result,
	); err != nil {
		t.Fatalf("WriteExplain: %v", err)
	}

	if strings.Contains(buffer.String(), "Pack policy") {
		t.Fatalf("empty Pack Policy must not render a section:\n%s", buffer.String())
	}
}
