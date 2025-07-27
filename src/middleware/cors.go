package middleware

import "github.com/gofiber/fiber/v2"

func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Set("Access-Control-Allow-Origin", origin)
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	}
}
