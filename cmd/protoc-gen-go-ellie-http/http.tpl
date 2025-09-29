{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}

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

        res, err := srv.{{.Name}}(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.HTTPStatusCodeFromError(err), hs.WrapHTTPResponse(res, err))
			ctx.Abort()
			return
		}

		hs.EncodeResponse(ctx, res, err)
    }
}
{{- end}}
