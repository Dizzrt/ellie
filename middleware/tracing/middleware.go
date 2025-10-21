package tracing

import (
	"context"
	"path"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type metadataCarrier struct {
	md metadata.MD
}

func (mc *metadataCarrier) Get(key string) string {
	values := mc.md.Get(key)
	if len(values) > 0 {
		return values[0]
	}

	return ""
}

func (mc *metadataCarrier) Set(key, value string) {
	mc.md.Set(key, value)
}

func (mc *metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc.md))
	for k := range mc.md {
		keys = append(keys, k)
	}
	return keys
}

const (
	httpTracerName       = "ellie/transport/http"
	grpcServerTracerName = "ellie/transport/grpc"
	grpcClientTracerName = "grpc-client"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			carrier := &metadataCarrier{md: md}
			ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
		}

		tracer := otel.Tracer(grpcServerTracerName)
		spanName := path.Base(info.FullMethod)

		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", path.Dir(info.FullMethod)[1:]),
			attribute.String("rpc.method", path.Base(info.FullMethod)),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.RecordError(err)
		}

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: add tracing for stream
		return handler(srv, ss)
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tracer := otel.Tracer(grpcClientTracerName)
		spanName := path.Base(method)

		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()

		span.SetAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", path.Dir(method)[1:]),
			attribute.String("rpc.method", path.Base(method)),
		)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		carrier := &metadataCarrier{md: md}
		otel.GetTextMapPropagator().Inject(ctx, carrier)
		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}

		return err
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// TODO: add tracing for stream
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// gin middleware for tracing
func Tracing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		rawCtx := req.Context()

		propagator := otel.GetTextMapPropagator()
		rawCtx = propagator.Extract(rawCtx, propagation.HeaderCarrier(req.Header))

		tracer := otel.Tracer(httpTracerName)
		rawCtx, span := tracer.Start(
			rawCtx,
			req.Method+" "+req.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(req.Method),
				semconv.HTTPRouteKey.String(req.URL.String()),
			),
		)
		defer span.End()

		ctx.Request = ctx.Request.WithContext(rawCtx)
		ctx.Next()

		span.SetAttributes(
			semconv.HTTPResponseStatusCodeKey.Int(ctx.Writer.Status()),
		)

		if len(ctx.Errors) > 0 {
			span.RecordError(ctx.Errors[0].Err)
		}
	}
}
