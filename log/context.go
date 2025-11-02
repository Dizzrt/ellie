package log

import (
	"context"

	"github.com/dizzrt/ellie/log/logid"
	"go.opentelemetry.io/otel/propagation"
)

type _CONTEXT_KEY_LOG_ID struct{}
type _CONTEXT_KEY_SPAN_ID struct{}
type _CONTEXT_KEY_TRACE_ID struct{}
type _CONTEXT_KEY_TRACEPARENT struct{}

var ctxKeyMappings = []struct {
	key  any
	name string
}{
	{_CONTEXT_KEY_LOG_ID{}, "log_id"},
	{_CONTEXT_KEY_SPAN_ID{}, "span_id"},
	{_CONTEXT_KEY_TRACE_ID{}, "trace_id"},
	{_CONTEXT_KEY_TRACEPARENT{}, "traceparent"},
}

func LogIDFromContext(ctx context.Context) string {
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

// extract log_id and traceparent and inject them into context
func ExtractFromTextMapCarrier(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	if logID := LogIDFromContext(ctx); logID == "" {
		if logID := carrier.Get("log.id"); logID != "" {
			// for grpc metadata
			ctx = WithLogID(ctx, logID)
		} else if logID := carrier.Get("X-Log-ID"); logID != "" {
			// for http header
			ctx = WithLogID(ctx, logID)
		} else {
			ctx = WithLogID(ctx, logid.Generate().String())
		}
	}

	traceparent := carrier.Get("traceparent")
	if traceparent != "" {
		ctx = context.WithValue(ctx, _CONTEXT_KEY_TRACEPARENT{}, traceparent)
	}

	return ctx
}
