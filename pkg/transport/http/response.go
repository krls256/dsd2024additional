package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/krls256/dsd2024additional/pkg/transport"

	"net/http"
)

func OK(ctx *fiber.Ctx, meta interface{}, data interface{}) error {
	r := transport.NewResponse(http.StatusOK, meta, nil, data)

	return ctx.Status(r.Status).JSON(r)
}

func OKRaw(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(http.StatusOK).JSON(data)
}

func File(ctx *fiber.Ctx, content []byte, filename string) error {
	ctx.Set("Content-Description", "File Transfer")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

	return ctx.Send(content)
}

func Redirect(ctx *fiber.Ctx, to string) error {
	return ctx.Redirect(to, http.StatusFound)
}

func RawBadRequest(ctx *fiber.Ctx, err interface{}) error {
	return ctx.Status(http.StatusBadRequest).JSON(err)
}

func BadRequest(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusBadRequest, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func Unauthorized(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusUnauthorized, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func PaymentRequired(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusPaymentRequired, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func Forbidden(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusForbidden, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func NotFound(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusNotFound, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func Conflict(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusConflict, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func Teapot(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusTeapot, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func ValidationFailed(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusUnprocessableEntity, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func TooEarly(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusTooEarly, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func TooManyRequests(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusTooManyRequests, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}

func ServerError(ctx *fiber.Ctx, meta interface{}, err ...error) error {
	r := transport.NewFormattedError(http.StatusInternalServerError, meta, nil, err)

	return ctx.Status(r.Status).JSON(r)
}
