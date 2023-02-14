package fiber_server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func loggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()
	zap.L().Info("request log",
		zap.Duration("duration", time.Since(start)),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("ip", c.IP()),
		zap.Int("status", c.Response().StatusCode()),
		zap.String("hostname", c.Hostname()),
		zap.ByteString("header", c.Request().Header.Header()),
	)

	return err
}
