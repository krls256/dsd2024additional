package di

import (
	"github.com/krls256/dsd2024additional/internal/chat/constants"
	"github.com/krls256/dsd2024additional/internal/chat/services"
	"github.com/krls256/dsd2024additional/internal/chat/ws"
	pkgConstants "github.com/krls256/dsd2024additional/pkg/constants"
	pkgDI "github.com/krls256/dsd2024additional/pkg/di"
	"github.com/krls256/dsd2024additional/pkg/redis"
	"github.com/krls256/dsd2024additional/pkg/transport/hub"
	"github.com/sarulabs/di/v2"
)

func Defs() []di.Def {
	return []di.Def{
		{
			Name: constants.ChatServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(pkgConstants.RedisName).(*redis.Client)
				h := ctn.Get(pkgConstants.HubName).(*hub.Hub)

				return services.NewChatService(h, conn), nil
			},
		},
		{
			Name: constants.ChatHandlerName,
			Tags: []di.Tag{{Name: pkgDI.WSHandlerTag}},
			Build: func(ctn di.Container) (interface{}, error) {
				chatService := ctn.Get(constants.ChatServiceName).(*services.ChatService)

				return ws.NewChatHandler(chatService), nil
			},
		},
	}
}
