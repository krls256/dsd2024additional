package ws

import (
	"github.com/krls256/dsd2024additional/internal/chat/entities"
	"github.com/krls256/dsd2024additional/internal/chat/services"
	"github.com/krls256/dsd2024additional/pkg/transport/websocket"
	"github.com/krls256/dsd2024additional/utils"
	"go.uber.org/zap"
)

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

type ChatHandler struct {
	chatService *services.ChatService
}

func (h *ChatHandler) Register(r *websocket.Router) {
	r.Accept("push", func(ctx *websocket.Context) {
		m, err := utils.ReMarshal[entities.Message](ctx.Payload())
		if err != nil {
			zap.S().Info(1)
			ctx.ResponseWriter().BadRequest("push", nil, err)

			return
		}

		if err := h.chatService.Push(ctx, m); err != nil {
			ctx.ResponseWriter().BadRequest("push", nil, err)

			return
		}
	})
}
