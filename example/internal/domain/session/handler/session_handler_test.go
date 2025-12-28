package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mickamy/gokitx/either"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mickamy/contest"
	authv1 "github.com/mickamy/contest/example/gen/auth/v1"
	"github.com/mickamy/contest/example/gen/github.com/mickamy/contest/example/gen/auth/v1/authv1connect"
	"github.com/mickamy/contest/example/internal/domain/session/handler"
	"github.com/mickamy/contest/example/internal/domain/session/usecase"
	"github.com/mickamy/contest/example/internal/domain/session/usecase/mock_usecase"
	"github.com/mickamy/contest/example/test/cerrors"
)

func TestSession_SignIn(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		uc        func(ctrl *gomock.Controller) usecase.CreateSession
		wantCode  int
		assertRes func(t *testing.T, res *authv1.SignInResponse)
		assertErr func(t *testing.T, err *connect.Error)
	}{
		{
			name: "success",
			uc: func(ctrl *gomock.Controller) usecase.CreateSession {
				return usecase.NewCreateSession()
			},
			wantCode: http.StatusOK,
			assertRes: func(t *testing.T, res *authv1.SignInResponse) {
				require.NotZero(t, res.Tokens)
				assert.NotZero(t, res.Tokens.Access)
				assert.NotZero(t, res.Tokens.Refresh)
			},
			assertErr: func(t *testing.T, err *connect.Error) {
				require.NoError(t, err)
			},
		},
		{
			name: "bad request error",
			uc: func(ctrl *gomock.Controller) usecase.CreateSession {
				uc := mock_usecase.NewMockCreateSession(ctrl)
				uc.EXPECT().Do(gomock.Any(), gomock.Any()).Return(usecase.CreateSessionOutput{}, errors.New("bad request"))
				return uc
			},
			wantCode: http.StatusBadRequest,
			assertErr: func(t *testing.T, err *connect.Error) {
				require.Error(t, err)
				cerrors.AssertCode(t, connect.CodeInvalidArgument, err)
				errDetails := cerrors.ExtractErrorDetails(t, err)
				assert.Equal(t, "bad request", errDetails.Message)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var out authv1.SignInResponse
			ct := contest.
				New(t, either.Right(authv1connect.NewSessionServiceHandler(handler.NewSession(tc.uc(ctrl))))).
				Procedure(authv1connect.SessionServiceSignInProcedure).
				In(&authv1.SignInRequest{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 12),
				}).
				Do().
				ExpectStatus(tc.wantCode)

			if tc.assertRes != nil {
				ct.Out(&out)
				tc.assertRes(t, &out)
			}
			if tc.assertErr != nil {
				tc.assertErr(t, ct.Err())
			}
		})
	}
}
