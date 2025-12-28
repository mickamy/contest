package response

import (
	"fmt"

	"connectrpc.com/connect"

	commonv1 "github.com/mickamy/contest/example/gen/common/v1"
)

type InternalError struct {
	Message         string
	underlyingError error
}

func NewInternalError(underlyingError error) *InternalError {
	return &InternalError{
		Message:         "An internal error has occurred.",
		underlyingError: underlyingError,
	}
}

func (e *InternalError) Error() string {
	return e.underlyingError.Error()
}

func (e *InternalError) WithMessage(message string) *InternalError {
	e.Message = message
	return e
}

func (e *InternalError) AsConnectError() *connect.Error {
	connErr := connect.NewError(connect.CodeInternal, e.underlyingError)

	proto := e.asProto()
	detail, err := connect.NewErrorDetail(proto)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create error detail: %w", err))
	}

	connErr.AddDetail(detail)
	return connErr
}

func (e *InternalError) asProto() *commonv1.ErrorDetails {
	return &commonv1.ErrorDetails{
		Message: e.Message,
	}
}
