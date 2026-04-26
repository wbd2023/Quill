package coverage

import "testing"

func TestCoverageIncludesEveryStyleHeading(t *testing.T) {
	document := loadDocument(t)
	covered := make(map[string]bool)

	report := loadCoverageReport(t)
	for _, entry := range report.Sections {
		covered[entry.Section] = true
	}

	for _, heading := range document.Headings {
		section := heading.Section
		if covered[section] {
			continue
		}

		t.Fatalf("STYLE.md section %q missing from coverage index", section)
	}
}
