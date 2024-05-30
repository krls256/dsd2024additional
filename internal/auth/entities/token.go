package entities

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`

	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (t Token) GetID() uuid.UUID {
	return t.ID
}

func (t Token) IDColumnName() string {
	return "id"
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
