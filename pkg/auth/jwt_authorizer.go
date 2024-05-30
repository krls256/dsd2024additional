package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWTAuthorizer struct {
	cfg *JWTConfig
}

func NewAuthorizer(cfg *JWTConfig) *JWTAuthorizer {
	return &JWTAuthorizer{
		cfg: cfg,
	}
}

func (p *JWTAuthorizer) Token(opts ...TokenOption) (*Token, error) {
	options := NewTokenOptions(opts...)
	expiresAt := time.Now().Add(options.Expiry)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        options.ID,
		Issuer:    p.cfg.Issuer,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	token, err := t.SignedString([]byte(p.cfg.Fingerprint))
	if err != nil {
		return nil, err
	}

	return &Token{
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

func (p *JWTAuthorizer) Refresh(opts ...TokenOption) (*JWT, error) {
	options := NewTokenOptions(opts...)

	secret := options.RefreshToken
	if len(options.Secret) > 0 {
		secret = options.Secret
	}

	if _, err := p.Inspect(secret); err != nil {
		return nil, err
	}

	access, err := p.Token(WithExpiry(p.cfg.AccessTokenLifetime), WithTokenID(options.ID))
	if err != nil {
		return nil, err
	}

	refresh, err := p.Token(WithExpiry(p.cfg.RefreshTokenLifetime), WithTokenID(options.ID))
	if err != nil {
		return nil, err
	}

	return &JWT{
		CreatedAt:    access.CreatedAt,
		ExpiresAt:    access.ExpiresAt,
		AccessToken:  access.Token,
		RefreshToken: refresh.Token,
	}, nil
}

func (p *JWTAuthorizer) Generate(opts ...GenerateOption) (*JWT, error) {
	options := NewGenerateOptions(opts...)

	access, err := p.Token(WithExpiry(p.cfg.AccessTokenLifetime), WithTokenID(options.ID))
	if err != nil {
		return nil, err
	}

	refresh, err := p.Token(WithExpiry(p.cfg.RefreshTokenLifetime))
	if err != nil {
		return nil, err
	}

	return &JWT{
		CreatedAt:    access.CreatedAt,
		ExpiresAt:    access.ExpiresAt,
		AccessToken:  access.Token,
		RefreshToken: refresh.Token,
	}, nil
}

func (p *JWTAuthorizer) Verify(token string) (string, error) {
	return p.Inspect(token)
}

func (p *JWTAuthorizer) Inspect(t string) (string, error) {
	token, err := p.parse(t)
	if token != nil && token.Valid {
		return token.Claims.(*jwt.RegisteredClaims).ID, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return "", jwt.ErrTokenExpired
		}
	}

	return "", jwt.ErrTokenUnverifiable
}

func (p *JWTAuthorizer) parse(t string) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(t, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}

		return []byte(p.cfg.Fingerprint), nil
	})

	return token, err
}
