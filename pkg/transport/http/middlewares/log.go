package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
)

type LoggerConfig struct {
	LogErrorsOnly bool
	LogResponse   bool
}

func log(ctx *fiber.Ctx, logResponse bool) {
	args := []interface{}{
		ctx.IP(), " | ", ctx.Response().StatusCode(), " | ", ctx.Method(), " | ", ctx.OriginalURL(),
		"\nREQUEST:\n", string(ctx.Body()),
	}

	if logResponse {
		args = append(args, "\nRESPONSE:\n", string(ctx.Response().Body()))
	}

	zap.S().Info(args...)
}

func RequestLoggerMiddleware(cfg LoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if cfg.LogErrorsOnly {
			if err != nil || c.Response().StatusCode() >= http.StatusBadRequest {
				log(c, cfg.LogResponse)
			}
		} else {
			log(c, cfg.LogResponse)
		}

		return err
	}
}
