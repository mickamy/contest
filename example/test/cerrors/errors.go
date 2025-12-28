package cerrors

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/mickamy/gokitx/either"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonv1 "github.com/mickamy/contest/example/gen/common/v1"
	"github.com/mickamy/contest/example/gen/google/rpc"
)

func AssertCode(t *testing.T, expected connect.Code, connErr *connect.Error) {
	t.Helper()

	assert.Equalf(t, expected, connect.CodeOf(connErr), "code=%s", connect.CodeOf(connErr).String())
}

func ExtractErrorDetails(t *testing.T, connErr *connect.Error) *commonv1.ErrorDetails {
	t.Helper()

	require.Len(t, connErr.Details(), 1)

	detail := either.Must(connErr.Details()[0].Value())
	if commonErr, ok := detail.(*commonv1.ErrorDetails); ok {
		return commonErr
	}

	t.Fatalf("unexpected error type: %T", detail)
	return nil
}

func AsBadRequest(t *testing.T, errDetails *commonv1.ErrorDetails) *rpc.BadRequest {
	t.Helper()

	payload, ok := errDetails.Payload.(*commonv1.ErrorDetails_BadRequest)
	if !ok {
		t.Fatalf("payload type %T is not commonv1.ErrorDetails_BadRequest", errDetails.Payload)
	}

	return payload.BadRequest
}
