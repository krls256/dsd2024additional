package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krls256/dsd2024additional/internal/auth/entities"
	"github.com/krls256/dsd2024additional/internal/auth/services"
	"github.com/krls256/dsd2024additional/pkg/auth"
	"github.com/krls256/dsd2024additional/pkg/errors"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
)

func NewAccountHandler(accountService *services.AccountService, errorHandler *errors.ErrorHTTPHandler,
	jwtFactory *auth.JWTMiddlewareFactory) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		errorHandler:   errorHandler,
		jwtFactory:     jwtFactory,
	}
}

type AccountHandler struct {
	accountService *services.AccountService
	errorHandler   *errors.ErrorHTTPHandler
	jwtFactory     *auth.JWTMiddlewareFactory
}

func (h *AccountHandler) Register(router fiber.Router) {
	account := router.Group("auth/account")

	account.Put("", h.create)
	account.Delete("", h.delete)
}

func (h *AccountHandler) create(ctx *fiber.Ctx) error {
	req := entities.CreateAccount{}
	if err := ctx.BodyParser(&req); err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	acc, err := h.accountService.Create(ctx.UserContext(), req)
	if err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	return http.OK(ctx, nil, acc)
}

func (h *AccountHandler) delete(ctx *fiber.Ctx) error {
	req := entities.DeleteAccountRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	if err := h.accountService.Delete(ctx.UserContext(), req); err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	return http.OK(ctx, nil, nil)
}
