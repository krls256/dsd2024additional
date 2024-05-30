package auth

import "time"

type JWT struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (t *JWT) Expired() bool {
	return t.ExpiresAt.Unix() < time.Now().Unix()
}
