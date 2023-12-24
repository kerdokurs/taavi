package middleware

import "github.com/gofiber/fiber/v2"

const htmxHeader = "Hx-Request"

func Htmx() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()

		_, isHtmx := headers[htmxHeader]
		c.Locals("htmx", isHtmx)

		return c.Next()
	}
}
