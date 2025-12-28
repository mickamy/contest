package contest

import (
	"fmt"

	"connectrpc.com/connect"
)

func sToConnectCode(s string) (connect.Code, error) {
	switch s {
	case "canceled":
		return connect.CodeCanceled, nil
	case "unknown":
		return connect.CodeUnknown, nil
	case "invalid_argument":
		return connect.CodeInvalidArgument, nil
	case "deadline_exceeded":
		return connect.CodeDeadlineExceeded, nil
	case "not_found":
		return connect.CodeNotFound, nil
	case "already_exists":
		return connect.CodeAlreadyExists, nil
	case "permission_denied":
		return connect.CodePermissionDenied, nil
	case "resource_exhausted":
		return connect.CodeResourceExhausted, nil
	case "failed_precondition":
		return connect.CodeFailedPrecondition, nil
	case "aborted":
		return connect.CodeAborted, nil
	case "out_of_range":
		return connect.CodeOutOfRange, nil
	case "unimplemented":
		return connect.CodeUnimplemented, nil
	case "internal":
		return connect.CodeInternal, nil
	case "unavailable":
		return connect.CodeUnavailable, nil
	case "data_loss":
		return connect.CodeDataLoss, nil
	case "unauthenticated":
		return connect.CodeUnauthenticated, nil
	}
	return connect.CodeUnknown, fmt.Errorf("unknown code: %s", s)
}
