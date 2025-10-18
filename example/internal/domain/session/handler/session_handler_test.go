package handler_test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mickamy/gokitx/either"

	"github.com/mickamy/connecttest"
	authv1 "github.com/mickamy/connecttest-example/gen/auth/v1"
	"github.com/mickamy/connecttest-example/gen/github.com/mickamy/connecttest-example/gen/auth/v1/authv1connect"
	"github.com/mickamy/connecttest-example/internal/domain/session/handler"
)

func TestSession_SignIn(t *testing.T) {
	t.Parallel()

	var out authv1.SignInResponse
	connecttest.
		New(t, either.Right(authv1connect.NewSessionServiceHandler(handler.NewSession()))).
		Procedure(authv1connect.SessionServiceSignInProcedure).
		In(&authv1.SignInRequest{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 12),
		}).
		Do().
		ExpectStatus(http.StatusOK).
		Out(&out)

	if out.Tokens == nil {
		t.Fatal("expected tokens to be set")
	}

	if out.Tokens.Access == nil {
		t.Fatal("expected access token to be set")
	}

	if out.Tokens.Refresh == nil {
		t.Fatal("expected refresh token to be set")
	}

	t.Logf("Access Token: %s", out.Tokens.Access.Value)
	t.Logf("Refresh Token: %s", out.Tokens.Refresh.Value)
}
