package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
)

type SessionRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Session, error)
	Set(ctx context.Context, session *entities.Session) error
	Delete(ctx context.Context, ids []uuid.UUID) error
}
