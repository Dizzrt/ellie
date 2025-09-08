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
	r.GET("/ping", _Ping_Ping_HTTP_Handler(srv))

	s.Handler = r
}

func _Ping_Ping_HTTP_Handler(srv PingHTTPServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req PingRequest

		data, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(nhttp.StatusBadRequest, err.Error())
			return
		}

		if len(data) > 0 {
			if err := ctx.ShouldBindJSON(&req); err != nil {
				ctx.JSON(nhttp.StatusBadRequest, err.Error())
				return
			}
		}

		res, err := srv.Ping(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(nhttp.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(nhttp.StatusOK, res)
	}
}
