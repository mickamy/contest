package response

import (
	"connectrpc.com/connect"
	"github.com/mickamy/gokitx/slices"

	commonv1 "github.com/mickamy/contest/example/gen/common/v1"
	"github.com/mickamy/contest/example/gen/google/rpc"
)

type BadRequestError struct {
	Err             error
	Message         string
	FieldViolations []FieldViolation
}

func NewBadRequestError(underlyingErr error, fieldViolations ...FieldViolation) *BadRequestError {
	return &BadRequestError{
		Err:             underlyingErr,
		Message:         underlyingErr.Error(),
		FieldViolations: fieldViolations,
	}
}

func (e *BadRequestError) WithMessage(message string) *BadRequestError {
	e.Message = message
	return e
}

func (e *BadRequestError) WithFieldViolation(field string, description ...string) *BadRequestError {
	e.FieldViolations = append(e.FieldViolations, FieldViolation{
		Field:        field,
		Descriptions: description,
	})
	return e
}

func (e *BadRequestError) AsConnectError() *connect.Error {
	connErr := connect.NewError(connect.CodeInvalidArgument, e.Err)
	proto, err := e.asProto()
	if err != nil {
		return NewInternalError(err).AsConnectError()
	}

	detail, err := connect.NewErrorDetail(proto)
	if err != nil {
		return NewInternalError(err).AsConnectError()
	}
	connErr.AddDetail(detail)

	return connErr
}

func (e *BadRequestError) asProto() (*commonv1.ErrorDetails, error) {
	rootErr := &commonv1.ErrorDetails{
		Message: e.Message,
	}
	var violations []*rpc.BadRequest_FieldViolation
	for _, violation := range e.FieldViolations {
		ps := violation.AsProto()
		violations = append(violations, ps...)
	}
	if violations != nil {
		badReq := &rpc.BadRequest{
			FieldViolations: violations,
		}

		rootErr.Payload = &commonv1.ErrorDetails_BadRequest{BadRequest: badReq}
	}

	return rootErr, nil
}

type FieldViolation struct {
	Field        string
	Descriptions []string
}

func (m FieldViolation) AsProto() []*rpc.BadRequest_FieldViolation {
	return slices.Map(m.Descriptions, func(description string) *rpc.BadRequest_FieldViolation {
		return &rpc.BadRequest_FieldViolation{
			Field:       m.Field,
			Description: description,
		}
	})
}
