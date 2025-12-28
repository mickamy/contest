package response_test

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/mickamy/gokitx/either"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonv1 "github.com/mickamy/contest/example/gen/common/v1"
	"github.com/mickamy/contest/example/internal/domain/common/dto/response"
)

func TestBadRequestError_AsConnectError(t *testing.T) {
	t.Parallel()

	// arrange
	err := response.NewBadRequestError(errors.New("test error")).
		WithMessage("test message").
		WithFieldViolation("test_field", "test description")

	// act
	connectErr := err.AsConnectError()

	// assert
	require.NotNil(t, connectErr)
	assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(connectErr), "code=%s", connect.CodeOf(connectErr).String())
	require.Equal(t, 1, len(connectErr.Details()))
	d := either.Must(connectErr.Details()[0].Value())
	errDetails, ok := d.(*commonv1.ErrorDetails)
	require.True(t, ok, "expected commonv1.ErrorDetails, got %T", d)

	payload, ok := errDetails.Payload.(*commonv1.ErrorDetails_BadRequest)
	require.Truef(t, ok, "expected BadRequest payload, got: %T", errDetails.Payload)

	badReq := payload.BadRequest
	assert.Equal(t, 1, len(badReq.FieldViolations))
	assert.Equal(t, "test_field", badReq.FieldViolations[0].Field)
	assert.Equal(t, "test description", badReq.FieldViolations[0].Description)
}
