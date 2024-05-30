package websocket

import (
	"fmt"
	"github.com/krls256/dsd2024additional/pkg/auth"
	"github.com/krls256/dsd2024additional/pkg/transport/hub"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"sync"
)

var wsUpgrader = websocket.Config{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
}

func NewServer(h *hub.Hub, jwtFactory *auth.JWTMiddlewareFactory, handlers ...Handler) *Server {
	s := &Server{
		h:          h,
		jwtFactory: jwtFactory,

		router:   NewRouter(),
		handlers: handlers,

		shutdown: make(chan struct{}),
	}

	for _, handler := range handlers {
		handler.Register(s.router)
	}

	for route := range s.router.routeMap {
		fmt.Printf("[WS route] %v\n", route)
	}

	return s
}

type Server struct {
	h          *hub.Hub
	jwtFactory *auth.JWTMiddlewareFactory

	router   *Router
	handlers []Handler

	shutdown chan struct{}
}

func (s *Server) Handler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID := s.jwtFactory.UnwrapCtx(ctx)

		return websocket.New(func(conn *websocket.Conn) {
			c := NewClient(conn, s.router, &sync.WaitGroup{})

			gorutinesCount := 2

			c.wg.Add(gorutinesCount)

			go c.writer()
			go c.reader(time.Minute * 30)

			s.h.Register(userID, c.rf.NewResponseWriter().AsyncWriter())
			defer s.h.Unregister(userID)

			c.wg.Wait()
		}, wsUpgrader)(ctx)
	}
}
