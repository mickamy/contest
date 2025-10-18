package handler_test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mickamy/gokitx/either"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/connecttest"
	authv1 "github.com/mickamy/connecttest-example/gen/auth/v1"
	"github.com/mickamy/connecttest-example/gen/github.com/mickamy/connecttest-example/gen/auth/v1/authv1connect"
	"github.com/mickamy/connecttest-example/internal/domain/session/handler"
	"github.com/mickamy/connecttest-example/internal/domain/session/usecase"
)

func TestSession_SignIn(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name      string
		uc        usecase.CreateSession
		wantCode  int
		assertRes func(t *testing.T, res *authv1.SignInResponse)
		assertErr func(t *testing.T, err error)
	}{
		{
			name:     "success",
			uc:       usecase.NewCreateSession(),
			wantCode: http.StatusOK,
			assertRes: func(t *testing.T, res *authv1.SignInResponse) {
				require.NotZero(t, res.Tokens)
				assert.NotZero(t, res.Tokens.Access)
				assert.NotZero(t, res.Tokens.Refresh)
			},
			assertErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var out authv1.SignInResponse
			connecttest.
				New(t, either.Right(authv1connect.NewSessionServiceHandler(handler.NewSession(tc.uc)))).
				Procedure(authv1connect.SessionServiceSignInProcedure).
				In(&authv1.SignInRequest{
					Email:    gofakeit.Email(),
					Password: gofakeit.Password(true, true, true, true, false, 12),
				}).
				Do().
				ExpectStatus(http.StatusOK).
				Out(&out)

			tc.assertRes(t, &out)
		})
	}
}
