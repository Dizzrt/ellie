package ping

import (
	context "context"

	nhttp "net/http"

	"github.com/Dizzrt/ellie/transport/http"
	"github.com/gin-gonic/gin"
)

type PingHTTPServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
}

func RegisterPingHTTPServer(s *http.Server, srv PingHTTPServer) {
	r := gin.Default()
	r.GET("/ping", _Ping_Ping_HTTP_Handler(s, srv))

	s.Handler = r
}

func _Ping_Ping_HTTP_Handler(s *http.Server, srv PingHTTPServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req PingRequest
		if err := s.BindPathParams(ctx.Request, &req); err != nil {
			ctx.JSON(nhttp.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		if err := s.BindQueryParams(ctx.Request, &req); err != nil {
			ctx.JSON(nhttp.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}

		res, err := srv.Ping(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(nhttp.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(nhttp.StatusOK, res)
	}
}
