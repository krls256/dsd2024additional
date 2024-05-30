package auth

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
	"github.com/samber/lo"
	"strings"
)

var (
	ErrPermissionDenied     = errors.New("permission denied")
	ErrAuthHeaderIsRequired = errors.New("auth header is required")
)

type SessionService interface {
	GetPermissions(ctx context.Context, id uuid.UUID) ([]string, error)
}

func NewJWTMiddlewareFactory(pm *PermissionManager, authorizer *JWTAuthorizer, cfg *JWTConfig, ss SessionService) *JWTMiddlewareFactory {
	return &JWTMiddlewareFactory{pm: pm, authorizer: authorizer, cfg: cfg, ss: ss}
}

type JWTMiddlewareFactory struct {
	pm         *PermissionManager
	authorizer *JWTAuthorizer
	cfg        *JWTConfig
	ss         SessionService
}

func (f *JWTMiddlewareFactory) GetToken(ctx *fiber.Ctx) string {
	headers := ctx.GetReqHeaders()[f.cfg.HeaderName]

	token := ""

	if len(headers) != 0 {
		token = headers[0]
	}

	return strings.ReplaceAll(token, f.cfg.HeaderScheme, "")
}

func (f *JWTMiddlewareFactory) Middleware(permission, description string) fiber.Handler {
	f.pm.Add(permission, description)

	return func(ctx *fiber.Ctx) error {
		t := f.GetToken(ctx)

		if t == "" {
			return http.Forbidden(ctx, nil, ErrAuthHeaderIsRequired)
		}

		jti, err := f.authorizer.Verify(t)
		if err != nil {
			return http.Forbidden(ctx, nil, err)
		}

		jtiUUID, err := uuid.Parse(jti)
		if err != nil {
			return http.Forbidden(ctx, nil, err)
		}

		permissions, err := f.ss.GetPermissions(ctx.UserContext(), jtiUUID)
		if err != nil {
			return http.Forbidden(ctx, nil, err)
		}

		if !lo.Contains(permissions, permission) {
			return http.Forbidden(ctx, nil, ErrPermissionDenied)
		}

		return ctx.Next()
	}
}
