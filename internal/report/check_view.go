package report

import "ciphera/tools/internal/style"

func NewCheckEntry(
	rule style.Rule,
	status style.CheckStatus,
	result style.ExecutionResult,
) (entry CheckEntry) {
	return CheckEntry{
		Rule:   NewRuleSummary(rule),
		Status: status,
		Result: result,
	}
}

func NewRuleSummary(rule style.Rule) (summary RuleSummary) {
	return RuleSummary{
		ID:             rule.ID,
		Name:           rule.Name,
		Group:          rule.Group,
		Enforcement:    rule.Enforcement,
		Scope:          rule.Scope,
		RequirementIDs: append([]string{}, rule.RequirementIDs...),
	}
}

func NewCheckView(result CheckResult) (view CheckView) {
	view = CheckView{
		Result:  result,
		Summary: result.Summary(),
		Groups:  make([]CheckGroup, 0),
	}

	for _, entry := range result.Entries {
		if len(view.Groups) == 0 ||
			entry.Rule.Group != view.Groups[len(view.Groups)-1].Group {
			view.Groups = append(view.Groups, CheckGroup{
				Group:   entry.Rule.Group,
				Entries: make([]CheckEntry, 0),
			})
		}

		lastIndex := len(view.Groups) - 1
		view.Groups[lastIndex].Entries = append(view.Groups[lastIndex].Entries, entry)
	}

	return view
}
