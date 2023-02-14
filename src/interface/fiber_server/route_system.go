package fiber_server

import "github.com/gofiber/fiber/v2"

func (f *FiberServer) addRouteSystem(base fiber.Router) {
	r := base.Group("/system")

	r.Get("/version", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte(f.config.AppVersion))
	})

	r.Get("/readiness", func(ctx *fiber.Ctx) error {
		err := f.useCase.HealthCheck(ctx.Context())
		if err != nil {
			return f.errorHandler(ctx, err)
		}

		return ctx.Send([]byte(OK))
	})

	r.Get("/liveliness", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte(OK))
	})

	r.Get("/liveness", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte(OK))
	})
}
