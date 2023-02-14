package fiber_server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (f *FiberServer) recover(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				f.errorHandler(c, fmt.Errorf("%v", r))
			} else {
				f.errorHandler(c, err)
			}
		}
	}()

	// Return err if exist, else move to next handler
	return c.Next()
}
