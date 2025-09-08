package errors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/status"
)

const (
	UnknownCode   = 500
	UnknownReason = ""
)

type Error struct {
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: %v", e.cause)
}

func New(code int, reason, message string) *Error {
	return &Error{}
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if ee := new(Error); errors.As(err, &ee) {
		return ee
	}

	grpcStatus, ok := status.FromError(err)
	if !ok {
		return New(UnknownCode, UnknownReason, err.Error())
	}

	ret := New(UnknownCode, UnknownReason, grpcStatus.Message())
	// for _, detail := range grpcStatus.Details() {
	// 	switch temp := detail.(type) {
	// 	case *errdetails.ErrorInfo:

	// 	}
	// }

	return ret
}
