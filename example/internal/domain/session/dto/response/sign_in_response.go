package response

import (
	authv1 "github.com/mickamy/contest/example/gen/auth/v1"
	"github.com/mickamy/contest/example/internal/domain/session/model"
)

func NewSignInResponse(m model.Session) *authv1.SignInResponse {
	return &authv1.SignInResponse{
		Tokens: &authv1.TokenSet{
			Access:  newToken(m.Tokens.Access),
			Refresh: newToken(m.Tokens.Refresh),
		},
	}
}
