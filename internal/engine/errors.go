package engine

import (
	"fmt"

	"github.com/wbd2023/Quill/internal/style"
)

func errUnknownScope(scope style.Scope) (err error) {
	return fmt.Errorf("unknown scope %q in style profile", scope)
}
