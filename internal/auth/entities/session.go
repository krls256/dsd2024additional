package entities

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *Session) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}
