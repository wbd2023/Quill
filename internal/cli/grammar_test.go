package cli

import "testing"

func TestCommandLineLookupCoversDeclaredCommands(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		args []string
	}{
		{name: "check", args: []string{"check"}},
		{name: "fix", args: []string{"fix"}},
		{name: "doctor", args: []string{"doctor"}},
		{name: "coverage", args: []string{"coverage"}},
		{name: "install", args: []string{"install"}},
		{name: "lock", args: []string{"lock"}},
		{name: "version", args: []string{"version"}},
		{name: "init", args: []string{"init"}},
		{name: "list", args: []string{"list", "packs"}},
		{name: "explain", args: []string{"explain", "rule:profile/enforcement-levels"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser, model, err := newParser()
			if err != nil {
				t.Fatalf("newParser: %v", err)
			}

			context, err := parser.Parse(testCase.args)
			if err != nil {
				t.Fatalf("Parse(%q): %v", testCase.args, err)
			}
			if model.lookup(context.Selected()) == nil {
				t.Fatalf("lookup(%q) = nil", testCase.name)
			}
		})
	}
}
