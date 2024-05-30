package http

import "github.com/gofiber/fiber/v2"

func NewWebHandler() WebHandler {
	return WebHandler{}
}

type WebHandler struct {
}

func (w WebHandler) Register(router fiber.Router) {
	router.Static("", "./internal/web/src/build")
}
