package handler

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	authv1 "github.com/mickamy/contest/example/gen/auth/v1"
	"github.com/mickamy/contest/example/gen/github.com/mickamy/contest/example/gen/auth/v1/authv1connect"
	commonResponse "github.com/mickamy/contest/example/internal/domain/common/dto/response"
	"github.com/mickamy/contest/example/internal/domain/session/dto/response"
	"github.com/mickamy/contest/example/internal/domain/session/usecase"
)

type Session struct {
	create usecase.CreateSession
}

func NewSession(
	create usecase.CreateSession,
) *Session {
	return &Session{
		create: create,
	}
}

func (h *Session) SignIn(ctx context.Context, req *connect.Request[authv1.SignInRequest]) (*connect.Response[authv1.SignInResponse], error) {
	out, err := h.create.Do(ctx, usecase.CreateSessionInput{
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
	})
	if err != nil {
		return nil, commonResponse.NewBadRequestError(err).AsConnectError()
	}

	res := connect.NewResponse(response.NewSignInResponse(out.Session))
	return res, nil
}

func (h *Session) Refresh(ctx context.Context, req *connect.Request[authv1.RefreshRequest]) (*connect.Response[authv1.RefreshResponse], error) {
	res := connect.NewResponse(&authv1.RefreshResponse{
		Tokens: &authv1.TokenSet{
			Access: &authv1.Token{
				Value:     uuid.New().String(),
				ExpiresAt: timestamppb.New(time.Now().Add(time.Hour)),
			},
			Refresh: &authv1.Token{
				Value:     uuid.New().String(),
				ExpiresAt: timestamppb.New(time.Now().Add(time.Hour * 24 * 7)),
			},
		},
	})
	return res, nil
}

func (h *Session) SignOut(ctx context.Context, req *connect.Request[authv1.SignOutRequest]) (*connect.Response[authv1.SignOutResponse], error) {
	return connect.NewResponse(&authv1.SignOutResponse{}), nil
}

var _ authv1connect.SessionServiceHandler = (*Session)(nil)
