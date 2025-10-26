package log

type logConstKey string

const (
	// Msg     logConstKey = "msg"
	LogID   logConstKey = "log_id"
	TraceID logConstKey = "trace_id"
	SpanID  logConstKey = "span_id"
)

var fromCtxKeys = []logConstKey{LogID, TraceID, SpanID}
