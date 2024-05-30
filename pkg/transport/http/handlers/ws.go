package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krls256/dsd2024additional/pkg/auth"
	"github.com/krls256/dsd2024additional/pkg/transport/websocket"
)

func NewWSHandler(wsServer *websocket.Server, jwtFactory *auth.JWTMiddlewareFactory) *WSHandler {
	return &WSHandler{wsServer: wsServer, jwtFactory: jwtFactory}
}

type WSHandler struct {
	wsServer   *websocket.Server
	jwtFactory *auth.JWTMiddlewareFactory
}

func (h *WSHandler) Register(router fiber.Router) {
	router.Get("ws", h.jwtFactory.Middleware(), h.wsServer.Handler())
}
