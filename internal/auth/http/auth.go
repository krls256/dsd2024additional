package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	"github.com/krls256/dsd2024additional/internal/auth/services"
	"github.com/krls256/dsd2024additional/pkg/auth"
	"github.com/krls256/dsd2024additional/pkg/errors"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
)

func NewAuthService(
	authService *services.AuthService,
	errorHandler *errors.ErrorHTTPHandler,
	jwtFactory *auth.JWTMiddlewareFactory) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		errorHandler: errorHandler,
		jwtFactory:   jwtFactory,
	}
}

type AuthHandler struct {
	authService  *services.AuthService
	errorHandler *errors.ErrorHTTPHandler
	jwtFactory   *auth.JWTMiddlewareFactory
}

func (h *AuthHandler) Register(router fiber.Router) {
	authGroup := router.Group("auth")
	authGroup.Post("login", h.login)

	authGroup.Post("refresh", h.refresh)
	authGroup.Post("logout", h.logout)
}

func (h *AuthHandler) login(ctx *fiber.Ctx) error {
	req := entities.LoginRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	tokens, err := h.authService.Login(ctx.UserContext(), req)
	if err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	return http.OK(ctx, nil, tokens)
}

func (h *AuthHandler) refresh(ctx *fiber.Ctx) error {
	t := h.jwtFactory.GetToken(ctx)

	tokens, err := h.authService.Refresh(ctx.UserContext(), t)
	if err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	return http.OK(ctx, nil, tokens)
}

func (h *AuthHandler) logout(ctx *fiber.Ctx) error {
	t := h.jwtFactory.GetToken(ctx)

	if err := h.authService.Logout(ctx.UserContext(), t); err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	return http.OK(ctx, nil, nil)
}
