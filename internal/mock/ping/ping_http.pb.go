package ping

import (
	context "context"

	"github.com/Dizzrt/ellie/transport/http"
	"github.com/gin-gonic/gin"
)

type PingHTTPServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
}

func RegisterPingHTTPServer(s *http.Server, srv PingHTTPServer) {
	r := s.Engine()

	r.GET("/ping", _Ping_Ping_HTTP_Handler(s, srv))
	r.POST("/hello/:name", _Ping_Hello_HTTP_handler(s, srv))

	s.Handler = r
}

func _Ping_Ping_HTTP_Handler(hs *http.Server, srv PingHTTPServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req PingRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, http.WrapHTTPResponse(nil, err))
			ctx.Abort()
			return
		}

		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, http.WrapHTTPResponse(nil, err))
			ctx.Abort()
			return
		}

		res, err := srv.Ping(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.HTTPStatusCodeFromError(err), http.WrapHTTPResponse(res, err))
			ctx.Abort()
			return
		}

		// ctx.JSON(http.StatusOK, http.WrapHTTPResponse(res, err))
		hs.EncodeResponse(ctx, res, err)
	}
}

func _Ping_Hello_HTTP_handler(hs *http.Server, srv PingHTTPServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req HelloRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, http.WrapHTTPResponse(nil, err))
			ctx.Abort()
			return
		}

		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, http.WrapHTTPResponse(nil, err))
			ctx.Abort()
			return
		}

		if ctx.Request.ContentLength > 0 {
			if err := ctx.ShouldBind(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, http.WrapHTTPResponse(nil, err))
				ctx.Abort()
				return
			}
		}

		res, err := srv.Hello(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.HTTPStatusCodeFromError(err), http.WrapHTTPResponse(res, err))
			ctx.Abort()
			return
		}

		// ctx.JSON(http.StatusOK, http.WrapHTTPResponse(res, err))
		hs.EncodeResponse(ctx, res, err)
	}
}
