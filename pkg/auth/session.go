package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pkgErrors "github.com/krls256/dsd2024additional/pkg/errors"
	"github.com/krls256/dsd2024additional/pkg/redis"
	"github.com/samber/lo"
	"time"
)

var ErrSessionsNotFound = pkgErrors.NewErrorWithCode("session not found", 104001)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *Session) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

const SessionKey = "auth.session_%v"

func NewSessionService(conn *redis.Client) *SessionService {
	return &SessionService{conn: conn}
}

type SessionService struct {
	conn *redis.Client
}

func (s *SessionService) Get(ctx context.Context, id uuid.UUID) (*Session, error) {
	key := fmt.Sprintf(SessionKey, id)

	res, err := s.conn.Get(ctx, key)
	if errors.Is(err, redis.Nil) {
		return nil, ErrSessionsNotFound
	}

	if err != nil {
		return nil, err
	}

	file := &Session{}

	return file, json.Unmarshal(res, file)
}

func (s *SessionService) Set(ctx context.Context, session *Session) error {
	key := fmt.Sprintf(SessionKey, session.ID)

	return s.conn.Set(ctx, key, session, time.Until(session.ExpiresAt))
}

func (s *SessionService) Delete(ctx context.Context, ids []uuid.UUID) error {
	keys := lo.Map(ids, func(item uuid.UUID, index int) string {
		return fmt.Sprintf(SessionKey, item)
	})

	return s.conn.Del(ctx, keys...)
}
