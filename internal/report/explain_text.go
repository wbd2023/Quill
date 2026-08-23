package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func writeExplainText(writer io.Writer, result ExplainResult) (err error) {
	rule := result.Rule

	if _, err = fmt.Fprintf(writer, "Rule: %s\n\n", rule.ID); err != nil {
		return err
	}

	rows := [][2]string{
		{"pack", fmt.Sprintf("%s (%s)", rule.Pack.ID, rule.Pack.Provenance)},
		{"name", rule.Name},
		{"group", rule.Group},
		{"enforcement", rule.Binding.Enforcement},
		{"scope", rule.Binding.Scope},
		{"requirements", strings.Join(rule.Binding.Requirements, ", ")},
	}
	rows = append(rows, executionRows("check", rule.Check)...)

	if rule.Fix == nil {
		rows = append(rows, [2]string{"fix", "none"})
	} else {
		rows = append(rows, executionRows("fix", *rule.Fix)...)
	}

	if err = writeFields(writer, rows); err != nil {
		return err
	}

	return writePackPolicy(writer, rule.Pack.Policy)
}

// writePackPolicy renders the relevant Pack Policy deterministically beneath the rule explanation.
// json.MarshalIndent emits map keys in sorted order, so the output is stable across runs and
// platforms. An empty or absent policy renders nothing.
func writePackPolicy(writer io.Writer, policy map[string]any) (err error) {
	if len(policy) == 0 {
		return nil
	}

	if _, err = fmt.Fprint(writer, "\nPack policy\n\n"); err != nil {
		return err
	}

	encoded, err := json.MarshalIndent(policy, "", indentPrefix)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(strings.TrimRight(string(encoded), "\n"), "\n") {
		if _, err = fmt.Fprintf(writer, "%s%s\n", indentPrefix, line); err != nil {
			return err
		}
	}

	return nil
}

func executionRows(label string, execution ExplainExecution) (rows [][2]string) {
	rows = append(rows, [2]string{label, execution.Category})
	if len(execution.Tools) > 0 {
		rows = append(rows, [2]string{label + " tools", strings.Join(execution.Tools, ", ")})
	}
	if execution.FileSet != "" {
		rows = append(rows, [2]string{label + " file set", execution.FileSet})
	}
	if execution.Language != "" {
		rows = append(rows, [2]string{label + " language", execution.Language})
	}

	return rows
}
