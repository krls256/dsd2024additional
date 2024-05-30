package entities

import (
	"github.com/google/uuid"
	"time"
)

func NewProfile(req UpsertProfileRequest) *Profile {
	p := &Profile{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	p.SetUpsertProfileRequest(req)

	return p
}

type Profile struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`

	Nickname string `json:"nickname"`
	AboutMe  string `json:"about_me"`
}

func (p *Profile) SetUpsertProfileRequest(req UpsertProfileRequest) {
	p.ID = req.ID
	p.Nickname = req.Nickname
	p.AboutMe = req.AboutMe
}

func (p *Profile) GetID() uuid.UUID {
	return p.ID
}

func (p *Profile) IDColumnName() string {
	return "id"
}

type UpsertProfileRequest struct {
	ID uuid.UUID `json:"id"`

	Nickname string `json:"nickname"`
	AboutMe  string `json:"about_me"`
}
