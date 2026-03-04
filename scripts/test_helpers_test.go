package scripts_test

import (
	"path/filepath"
	"runtime"
)

/* ------------------------------------------- Helpers ------------------------------------------ */

func currentScriptsDirectory() (directory string) {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}
