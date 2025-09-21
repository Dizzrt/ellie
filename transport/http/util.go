package http

import (
	"net/http"
	"slices"

	"github.com/Dizzrt/ellie/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validHTTPCode = []int{
	100, 101, 102, 103,
	200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
	300, 301, 302, 303, 304, 305, 306, 307, 308,
	400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410,
	411, 412, 413, 414, 415, 416, 417, 418, 421, 422, 423,
	424, 425, 426, 428, 429, 431, 451,
	500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511,
}

func HTTPStatusCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.OK:
			return http.StatusOK
		case codes.InvalidArgument:
			return http.StatusBadRequest
		case codes.NotFound:
			return http.StatusNotFound
		case codes.AlreadyExists:
			return http.StatusConflict
		case codes.PermissionDenied:
			return http.StatusForbidden
		case codes.Unauthenticated:
			return http.StatusUnauthorized
		case codes.ResourceExhausted:
			return http.StatusTooManyRequests
		case codes.FailedPrecondition:
			return http.StatusPreconditionFailed
		case codes.Aborted:
			return http.StatusConflict
		case codes.OutOfRange:
			return http.StatusBadRequest
		case codes.Unimplemented:
			return http.StatusNotImplemented
		case codes.Internal:
			return http.StatusInternalServerError
		case codes.Unavailable:
			return http.StatusServiceUnavailable
		case codes.DataLoss:
			return http.StatusInternalServerError
		default:
			return http.StatusInternalServerError
		}
	}

	if ee, ok := err.(*errors.Error); ok && slices.Contains(validHTTPCode, int(ee.Code)) {
		return int(ee.Code)
	}

	return http.StatusOK
}

func WrapHTTPResponse(code int, message string, data any, err error) gin.H {
	if err != nil {
		if ee, ok := err.(*errors.Error); ok {
			// ellie error
			code = int(ee.Code)
			message = ee.Message
		} else if st, ok := status.FromError(err); ok {
			// grpc error
			code = int(st.Code())
			message = st.Message()
		} else {
			// unknown error type
			code = -1
			message = err.Error()
		}
	}

	return gin.H{
		"data":    data,
		"status":  code,
		"message": message,
	}
}
