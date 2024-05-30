package auth

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
	"strings"
)

var (
	ErrPermissionDenied     = errors.New("permission denied")
	ErrAuthHeaderIsRequired = errors.New("auth header is required")
)

func NewJWTMiddlewareFactory(authorizer *JWTAuthorizer, cfg *JWTConfig) *JWTMiddlewareFactory {
	return &JWTMiddlewareFactory{authorizer: authorizer, cfg: cfg}
}

type JWTMiddlewareFactory struct {
	authorizer *JWTAuthorizer
	cfg        *JWTConfig
}

func (f *JWTMiddlewareFactory) GetToken(ctx *fiber.Ctx) string {
	headers := ctx.GetReqHeaders()[f.cfg.HeaderName]

	token := ""

	if len(headers) != 0 {
		token = headers[0]
	}

	return strings.ReplaceAll(token, f.cfg.HeaderScheme, "")
}

func (f *JWTMiddlewareFactory) WrapCtx(ctx *fiber.Ctx, id uuid.UUID) {
	ctx.SetUserContext(context.WithValue(ctx.UserContext(), "id", id))
}

func (f *JWTMiddlewareFactory) UnwrapCtx(ctx *fiber.Ctx) uuid.UUID {
	return ctx.UserContext().Value("id").(uuid.UUID)
}

func (f *JWTMiddlewareFactory) Middleware() fiber.Handler {
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

		f.WrapCtx(ctx, jtiUUID)

		return ctx.Next()
	}
}
