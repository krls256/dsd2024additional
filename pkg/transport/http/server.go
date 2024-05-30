package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Server struct {
	cfg      Config
	fiberApp *fiber.App
}

func NewServer(ctx context.Context, appName string, cfg Config, handlers []Handler, middlewares []fiber.Handler) *Server {
	s := &Server{
		cfg: cfg,
		fiberApp: fiber.New(fiber.Config{
			StrictRouting:     false,
			CaseSensitive:     false,
			EnablePrintRoutes: true,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			AppName:           appName,
		})}

	s.fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			zap.S().Error("http exception", e)
		},
	}))

	s.fiberApp.Use(cors.New(cors.ConfigDefault))

	for _, h := range middlewares {
		s.fiberApp.Use(h)
	}

	for _, h := range handlers {
		h.Register(s.fiberApp)
	}

	return s
}

func (s *Server) AsyncRun() {
	go func() {
		if err := s.Run(); err != nil {
			zap.S().Error(err)
		}
	}()
}

func (s *Server) Run() error {
	if !s.cfg.Silent {
		zap.S().Info("Running server...")
	}

	return s.fiberApp.Listen(s.cfg.DNS())
}

func (s *Server) Shutdown(ctx context.Context) error {
	if !s.cfg.Silent {
		zap.S().Info("Shutdown server...")
	}

	defer func() {
		if !s.cfg.Silent {
			zap.S().Info("Server successfully stopped.")
		}
	}()

	if err := s.fiberApp.ShutdownWithContext(ctx); err != nil {
		if !s.cfg.Silent {
			zap.S().Error("Server forced to shutdown:", err)
		}

		return err
	}

	return nil
}
