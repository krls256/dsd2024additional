package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krls256/dsd2024additional/internal/profile/entities"
	"github.com/krls256/dsd2024additional/internal/profile/services"
	"github.com/krls256/dsd2024additional/pkg/auth"
	"github.com/krls256/dsd2024additional/pkg/errors"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
)

func NewProfileHandler(profileService *services.ProfileService,
	errorHandler *errors.ErrorHTTPHandler,
	jwtFactory *auth.JWTMiddlewareFactory) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		errorHandler:   errorHandler,
		jwtFactory:     jwtFactory,
	}
}

type ProfileHandler struct {
	profileService *services.ProfileService
	errorHandler   *errors.ErrorHTTPHandler
	jwtFactory     *auth.JWTMiddlewareFactory
}

func (h *ProfileHandler) Register(router fiber.Router) {
	profile := router.Group("profile")

	profile.Get("get", h.jwtFactory.Middleware(), h.get)
	profile.Post("upsert", h.jwtFactory.Middleware(), h.upsert)
}

func (h *ProfileHandler) upsert(ctx *fiber.Ctx) error {
	id := h.jwtFactory.UnwrapCtx(ctx)
	req := entities.UpsertProfileRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	req.ID = id

	p, err := h.profileService.Upsert(ctx.UserContext(), req)
	if err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	return http.OK(ctx, nil, p)
}

func (h *ProfileHandler) get(ctx *fiber.Ctx) error {
	req := entities.GetProfileRequest{
		ID: h.jwtFactory.UnwrapCtx(ctx),
	}

	p, err := h.profileService.Get(ctx.UserContext(), req)
	if err != nil {
		return h.errorHandler.HandleError(ctx, err)
	}

	return http.OK(ctx, nil, p)
}
