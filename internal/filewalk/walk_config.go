package filewalk

// WalkConfig carries the filesystem traversal policy that filewalk needs: which directories to
// exclude and how to detect generated files.
type WalkConfig struct {
	ExcludedDirectories []string
	GeneratedMarker     string
}
