package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

var _ LogWriter = (*stdLogger)(nil)

type stdLogger struct {
	w         io.Writer
	isDiscard bool
	mu        sync.Mutex
	pool      *sync.Pool
}

func NewStdLogger(w io.Writer) LogWriter {
	return &stdLogger{
		w:         w,
		isDiscard: w == io.Discard,
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (logger *stdLogger) Write(level Level, keyvals ...any) error {
	if logger.isDiscard || len(keyvals) == 0 {
		return nil
	}

	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	buf := logger.pool.Get().(*bytes.Buffer)
	defer logger.pool.Put(buf)

	buf.WriteString(level.String())
	for i := 0; i < len(keyvals); i += 2 {
		fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}

	buf.WriteByte('\n')
	defer buf.Reset()

	logger.mu.Lock()
	defer logger.mu.Unlock()
	_, err := logger.w.Write(buf.Bytes())

	return err
}

func (logger *stdLogger) Close() error {
	return nil
}
