package installer

import (
	"io"
)

// swappingReader invokes swap the first time it is read, then streams data. It deterministically
// injects a parent-directory swap mid-write so the commit cannot rely on a path validated earlier.
// It is shared by the Unix and Windows parent-swap regressions.
type swappingReader struct {
	swap    func() error
	data    []byte
	swapped bool
	offset  int
}

func (r *swappingReader) Read(p []byte) (n int, err error) {
	if !r.swapped {
		r.swapped = true
		if err = r.swap(); err != nil {
			return 0, err
		}
	}
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.offset:])
	r.offset += n

	return n, nil
}
