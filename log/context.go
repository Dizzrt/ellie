package log

import "context"

type _CONTEXT_KEY_LOG_ID struct{}
type _CONTEXT_KEY_SPAN_ID struct{}
type _CONTEXT_KEY_TRACE_ID struct{}

var ctxKeyMappings = []struct {
	key  any
	name string
}{
	{_CONTEXT_KEY_LOG_ID{}, "log_id"},
	{_CONTEXT_KEY_SPAN_ID{}, "span_id"},
	{_CONTEXT_KEY_TRACE_ID{}, "trace_id"},
}

func LogID(ctx context.Context) string {
	if logID, ok := ctx.Value(_CONTEXT_KEY_LOG_ID{}).(string); ok {
		return logID
	}

	return ""
}

func WithLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, _CONTEXT_KEY_LOG_ID{}, logID)
}

func SpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value(_CONTEXT_KEY_SPAN_ID{}).(string); ok {
		return spanID
	}

	return ""
}

func WithSpanID(ctx context.Context, spanID string) context.Context {
	if spanID == "" {
		return ctx
	}

	return context.WithValue(ctx, _CONTEXT_KEY_SPAN_ID{}, spanID)
}

func TraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(_CONTEXT_KEY_TRACE_ID{}).(string); ok {
		return traceID
	}

	return ""
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}

	return context.WithValue(ctx, _CONTEXT_KEY_TRACE_ID{}, traceID)
}
