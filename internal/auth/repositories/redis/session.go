package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	"github.com/krls256/dsd2024additional/internal/auth/errs"
	"github.com/krls256/dsd2024additional/pkg/redis"
	"github.com/samber/lo"
	"time"
)

const SessionKey = "auth.session_%v"

func NewSessionRepository(conn *redis.Client) *SessionRepository {
	return &SessionRepository{conn: conn}
}

type SessionRepository struct {
	conn *redis.Client
}

func (r *SessionRepository) Get(ctx context.Context, id uuid.UUID) (*entities.Session, error) {
	key := fmt.Sprintf(SessionKey, id)

	res, err := r.conn.Get(ctx, key)
	if errors.Is(err, redis.Nil) {
		return nil, errs.ErrSessionsNotFound
	}

	if err != nil {
		return nil, err
	}

	file := &entities.Session{}

	return file, json.Unmarshal(res, file)
}

func (r *SessionRepository) Set(ctx context.Context, session *entities.Session) error {
	key := fmt.Sprintf(SessionKey, session.ID)

	return r.conn.Set(ctx, key, session, time.Until(session.ExpiresAt))
}

func (r *SessionRepository) Delete(ctx context.Context, ids []uuid.UUID) error {
	keys := lo.Map(ids, func(item uuid.UUID, index int) string {
		return fmt.Sprintf(SessionKey, item)
	})

	return r.conn.Del(ctx, keys...)
}
