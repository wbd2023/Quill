package style_test

import "testing"

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStyleRegistryDefaultRows(t *testing.T) {
	t.Parallel()

	rows, err := loadStyleRegistryRows("")
	if err != nil {
		t.Fatalf("load default style registry rows: %v", err)
	}

	expectedRows := expectedStyleRegistryRows()
	if len(rows) != len(expectedRows) {
		t.Fatalf("row count mismatch: got %d, want %d", len(rows), len(expectedRows))
	}

	for index := range expectedRows {
		if rows[index] != expectedRows[index] {
			t.Fatalf(
				"row %d mismatch:\n got:  %#v\n want: %#v",
				index,
				rows[index],
				expectedRows[index],
			)
		}
	}

	seen := map[styleRegistryRow]struct{}{}
	for _, row := range rows {
		if _, exists := seen[row]; exists {
			t.Fatalf("duplicate registry row detected: %#v", row)
		}
		seen[row] = struct{}{}
	}
}
