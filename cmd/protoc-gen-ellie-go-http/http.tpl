{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}

const TRACER_NAME = "{{$.PackagePath}}"

{{- range .MethodSets}}
const Operation{{$svrType}}{{.OriginalName}} = "/{{$svrName}}/{{.OriginalName}}"
{{- end}}

type {{.ServiceType}}HTTPServer interface {
{{- range .MethodSets}}
    {{- if ne .Comment ""}}
    {{.Comment}}
    {{- end}}
    {{.Name}}(context.Context, *{{.Request}}) (*{{.Response}}, error)
{{- end}}
}

func Register{{.ServiceType}}HTTPServer(hs *http.Server, srv {{.ServiceType}}HTTPServer) {
    r := hs.Engine()

    {{- range .Methods}}
    r.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}_{{.Num}}_HTTP_Handler(hs, srv))
    {{- end}}
}

{{- range .Methods}}
func _{{$svrType}}_{{.Name}}_{{.Num}}_HTTP_Handler(hs *http.Server, srv {{$svrType}}HTTPServer) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        var req {{.Request}}
        if err := ginx.DecodeRequest(ctx, &req); err != nil {
			ctx.JSON(http.StatusBadRequest, hs.WrapHTTPResponse(nil, err))
			ctx.Abort()
			return
		}

        greq := ctx.Request
		rctx := greq.Context()
		rctx = log.ExtractFromTextMapCarrier(rctx, propagation.HeaderCarrier(greq.Header))
		attributes := []attribute.KeyValue{
			v1_21_0.HTTPRequestMethodKey.String(greq.Method),
			v1_21_0.HTTPRouteKey.String(greq.URL.String()),
			attribute.String("log.id", log.LogIDFromContext(rctx)),
		}

        tracer := otel.Tracer(TRACER_NAME)
		rctx, span := tracer.Start(rctx, "_{{$svrType}}_{{.Name}}_{{.Num}}_HTTP_Handler",
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attributes...),
		)
		defer span.End()

        sctx := span.SpanContext()
		rctx = log.WithTraceID(rctx, sctx.TraceID().String())
		rctx = log.WithSpanID(rctx, sctx.SpanID().String())

        res, err := srv.{{.Name}}(rctx, &req)
        ctx.Request = ctx.Request.WithContext(rctx)
		if err != nil {
			ctx.JSON(http.HTTPStatusCodeFromError(err), hs.WrapHTTPResponse(res, err))
			ctx.Abort()
			return
		}

		hs.EncodeResponse(ctx, res, err)
    }
}
{{- end}}
