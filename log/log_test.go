package log

import "testing"

func TestStdLogger(t *testing.T) {
	logger, err := NewStdLoggerWriter("logs/log")
	if err != nil {
		t.Fatal(err)
	}

	logger.Write(LevelInfo, "hello")
	logger.(*stdLoggerWriter).Sync()
}
