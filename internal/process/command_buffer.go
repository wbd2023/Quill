package process

import (
	"strings"
	"sync"
)

// boundedBuffer is a strings.Builder wrapper that caps the total bytes written and records whether
// the limit was reached. It is safe for concurrent writers: os/exec copies stdout and stderr from
// independent goroutines, and the combined stream receives writes from both, so every mutation is
// guarded. Write always reports the full requested write count so producers are not short-circuited
// by truncation.
type boundedBuffer struct {
	mu        sync.Mutex
	builder   strings.Builder
	limit     int64
	written   int64
	truncated bool
}

func newBoundedBuffer(limit int64) (buffer *boundedBuffer) {
	return &boundedBuffer{limit: limit}
}

func (buffer *boundedBuffer) Write(data []byte) (count int, err error) {
	count = len(data)

	buffer.mu.Lock()
	defer buffer.mu.Unlock()

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

func (buffer *boundedBuffer) String() (output string) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()

	return buffer.builder.String()
}

// streamSink writes every byte to both the per-stream buffer and the combined buffer, so callers
// buffer is shared between the stdout and stderr sinks and is therefore reached by two goroutines;
// boundedBuffer serialises those writes.
type streamSink struct {
	stream   *boundedBuffer
	combined *boundedBuffer
}

func (sink *streamSink) Write(data []byte) (count int, err error) {
	_, _ = sink.stream.Write(data)
	_, _ = sink.combined.Write(data)
	return len(data), nil
}

func newOutputBuffers(limit int64) (stdout, stderr, combined *boundedBuffer) {
	return newBoundedBuffer(limit), newBoundedBuffer(limit), newBoundedBuffer(limit)
}
