package response

import (
	"connectrpc.com/connect"

	commonv1 "github.com/mickamy/contest/example/gen/common/v1"
)

type UnauthenticatedError struct {
	Err     error
	Message string
}

func NewUnauthenticatedError(underlyingErr error) *UnauthenticatedError {
	return &UnauthenticatedError{
		Err: underlyingErr,
	}
}

func (e *UnauthenticatedError) WithMessage(message string) *UnauthenticatedError {
	e.Message = message
	return e
}

func (e *UnauthenticatedError) AsConnectError() *connect.Error {
	connErr := connect.NewError(connect.CodeUnauthenticated, e.Err)
	proto := e.asProto()
	detail, err := connect.NewErrorDetail(proto)
	if err != nil {
		return NewInternalError(err).AsConnectError()
	}
	connErr.AddDetail(detail)
	return connErr
}

func (e *UnauthenticatedError) asProto() *commonv1.ErrorDetails {
	return &commonv1.ErrorDetails{
		Message: e.Message,
	}
}
