package response

import (
	"connectrpc.com/connect"

	commonv1 "github.com/mickamy/contest/example/gen/common/v1"
)

type NotFoundError struct {
	Err             error
	Message         string
	FieldViolations []FieldViolation
}

func NewNotFoundError(underlyingErr error) *NotFoundError {
	return &NotFoundError{Err: underlyingErr}
}

func (e *NotFoundError) WithMessage(message string) *NotFoundError {
	e.Message = message
	return e
}

func (e *NotFoundError) AsConnectError() *connect.Error {
	connErr := connect.NewError(connect.CodeNotFound, e.Err)
	proto := e.asProto()
	detail, err := connect.NewErrorDetail(proto)
	if err != nil {
		return NewInternalError(err).AsConnectError()
	}
	connErr.AddDetail(detail)
	return connErr
}

func (e *NotFoundError) asProto() *commonv1.ErrorDetails {
	return &commonv1.ErrorDetails{Message: e.Message}
}
