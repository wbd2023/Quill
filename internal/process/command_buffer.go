package process

import "strings"

// limitedBuffer is a strings.Builder wrapper that caps the total bytes written and tracks whether
// output was truncated.
type limitedBuffer struct {
	builder   strings.Builder
	limit     int64
	written   int64
	truncated bool
}

func (buffer *limitedBuffer) Write(data []byte) (count int, err error) {
	count = len(data)
	remaining := buffer.limit - buffer.written
	if remaining <= 0 {
		buffer.truncated = true
		return count, nil
	}

	if int64(len(data)) > remaining {
		data = data[:int(remaining)]
		buffer.truncated = true
	}

	buffer.written += int64(len(data))
	_, _ = buffer.builder.Write(data)
	return count, nil
}

func (buffer *limitedBuffer) String() (output string) {
	return buffer.builder.String()
}
