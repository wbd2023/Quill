package styleguide

import "testing"

func TestSourceFilePositionAt(t *testing.T) {
	file := newSourceFile("STYLE.md", []byte("ab\ncde\n\nf"))

	cases := []struct {
		name     string
		offset   int
		expected position
	}{
		{
			name:     "unknown offset",
			offset:   -1,
			expected: position{},
		},
		{
			name:     "first line",
			offset:   1,
			expected: position{line: 1, column: 2},
		},
		{
			name:     "second line",
			offset:   4,
			expected: position{line: 2, column: 2},
		},
		{
			name:     "blank line",
			offset:   7,
			expected: position{line: 3, column: 1},
		},
		{
			name:     "last line",
			offset:   8,
			expected: position{line: 4, column: 1},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := file.positionAt(test.offset)
			if got != test.expected {
				t.Fatalf("unexpected position\nexpected: %#v\nactual:   %#v", test.expected, got)
			}
		})
	}
}
