package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mickamy/contest/example/internal/domain/session/model"
)

type CreateSessionInput struct {
	Email    string
	Password string
}

type CreateSessionOutput struct {
	Session model.Session
}

//go:generate go tool mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateSession interface {
	Do(ctx context.Context, input CreateSessionInput) (CreateSessionOutput, error)
}

type createSession struct {
}

func NewCreateSession() CreateSession {
	return &createSession{}
}

func (uc *createSession) Do(ctx context.Context, input CreateSessionInput) (CreateSessionOutput, error) {
	userID := uuid.NewString()
	session, err := model.NewSession(userID)
	if err != nil {
		return CreateSessionOutput{}, fmt.Errorf("failed to create session: %w", err)
	}
	return CreateSessionOutput{Session: session}, nil
}
