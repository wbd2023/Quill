package engine

import (
	"fmt"

	"github.com/wbd2023/quill/internal/style"
)

// ArgumentError identifies a repository-dependent input that is syntactically valid but has no
// valid meaning for the loaded profile. The CLI maps it to its stable invalid_argument protocol
// code and usage exit status.
type ArgumentError struct {
	cause error
}

func (err *ArgumentError) Error() string {
	return err.cause.Error()
}

func (err *ArgumentError) Unwrap() error {
	return err.cause
}

func newArgumentError(format string, arguments ...any) (err error) {
	return &ArgumentError{cause: fmt.Errorf(format, arguments...)}
}

func errUnknownScope(scope style.Scope) (err error) {
	return newArgumentError("unknown scope %q in style profile", scope)
}
