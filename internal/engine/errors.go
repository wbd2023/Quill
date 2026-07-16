package engine

import (
	"fmt"

	"ciphera/tools/internal/style"
)

func errUnknownScope(scope style.Scope) (err error) {
	return fmt.Errorf("unknown scope %q in style profile", scope)
}
