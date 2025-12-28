package response

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	authv1 "github.com/mickamy/contest/example/gen/auth/v1"
	"github.com/mickamy/contest/example/internal/lib/jwt"
)

func newToken(m jwt.Token) *authv1.Token {
	return &authv1.Token{
		Value:     m.Value,
		ExpiresAt: timestamppb.New(m.ExpiresAt),
	}
}
