package runtime

import "testing"

func TestLayoutPathPrependsRepoToolBins(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")

	layout := LayoutForToolsDir("/repo/tools")
	if layout.StateDir != "/repo/.cache/style" {
		t.Fatalf("StateDir = %q, want %q", layout.StateDir, "/repo/.cache/style")
	}

	expected := layout.ToolBinDir + ":" + layout.NodeBinDir + ":/usr/bin"
	if actual := layout.SearchPath(); actual != expected {
		t.Fatalf("PATH = %q, want %q", actual, expected)
	}
}
