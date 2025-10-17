package handler

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	v1 "github.com/mickamy/connecttest-example/gen/auth/v1"
	"github.com/mickamy/connecttest-example/gen/github.com/mickamy/connecttest-example/gen/auth/v1/authv1connect"
)

type Session struct {
}

func NewSession() *Session {
	return &Session{}
}

func (h Session) SignIn(ctx context.Context, req *connect.Request[v1.SignInRequest]) (*connect.Response[v1.SignInResponse], error) {
	res := connect.NewResponse(&v1.SignInResponse{
		UserId: uuid.New().String(),
		Tokens: &v1.TokenSet{
			Access: &v1.Token{
				Value:     uuid.New().String(),
				ExpiresAt: timestamppb.New(time.Now().Add(time.Hour)),
			},
			Refresh: &v1.Token{
				Value:     uuid.New().String(),
				ExpiresAt: timestamppb.New(time.Now().Add(time.Hour * 24 * 7)),
			},
		},
	})
	return res, nil
}

func (h Session) Refresh(ctx context.Context, req *connect.Request[v1.RefreshRequest]) (*connect.Response[v1.RefreshResponse], error) {
	res := connect.NewResponse(&v1.RefreshResponse{
		Tokens: &v1.TokenSet{
			Access: &v1.Token{
				Value:     uuid.New().String(),
				ExpiresAt: timestamppb.New(time.Now().Add(time.Hour)),
			},
			Refresh: &v1.Token{
				Value:     uuid.New().String(),
				ExpiresAt: timestamppb.New(time.Now().Add(time.Hour * 24 * 7)),
			},
		},
	})
	return res, nil
}

func (h Session) SignOut(ctx context.Context, req *connect.Request[v1.SignOutRequest]) (*connect.Response[v1.SignOutResponse], error) {
	return connect.NewResponse(&v1.SignOutResponse{}), nil
}

var _ authv1connect.SessionServiceHandler = (*Session)(nil)
