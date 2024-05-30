package services

import (
	"context"
	"encoding/json"
	"github.com/krls256/dsd2024additional/internal/chat/entities"
	"github.com/krls256/dsd2024additional/pkg/redis"
	"github.com/krls256/dsd2024additional/pkg/transport/hub"
	"go.uber.org/zap"
	"time"
)

const ChatChannelName = "chat_channel"
const PushAction = "push"

func NewChatService(h *hub.Hub, conn *redis.Client) *ChatService {
	s := &ChatService{
		h:    h,
		conn: conn,
	}

	go s.run(context.Background())

	return s
}

type ChatService struct {
	h    *hub.Hub
	conn *redis.Client
}

func (s *ChatService) Push(ctx context.Context, message entities.Message) error {
	return s.conn.Publish(ctx, ChatChannelName, message)
}

func (s *ChatService) run(ctx context.Context) {
	ch := s.conn.Subscribe(ctx, ChatChannelName).Channel()

	for {
		message, ok := <-ch
		if !ok {
			time.Sleep(time.Second)

			ch = s.conn.Subscribe(context.Background(), ChatChannelName).Channel()
		}

		m := entities.Message{}

		if err := json.Unmarshal([]byte(message.Payload), &m); err != nil {
			zap.S().Error(err)

			continue
		}

		s.h.Broadcast(PushAction, m)
	}
}
