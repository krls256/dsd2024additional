package auth

import "time"

type JWTConfig struct {
	HeaderName           string
	QueryName            string
	HeaderScheme         string
	Fingerprint          string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	Issuer               string
}
