package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	"github.com/krls256/dsd2024additional/internal/auth/repositories"
)

func NewSessionService(sessionRepository repositories.SessionRepository) *SessionService {
	return &SessionService{sessionRepository: sessionRepository}
}

type SessionService struct {
	sessionRepository repositories.SessionRepository
}

func (s *SessionService) Get(ctx context.Context, id uuid.UUID) (*entities.Session, error) {
	return s.sessionRepository.Get(ctx, id)
}

func (s *SessionService) Set(ctx context.Context, session *entities.Session) error {
	return s.sessionRepository.Set(ctx, session)
}

func (s *SessionService) Delete(ctx context.Context, ids []uuid.UUID) error {
	return s.sessionRepository.Delete(ctx, ids)
}
