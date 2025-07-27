package middleware

import (
	"encoding/base64"
	"pg_bridge_go/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// BasicAuthMiddleware is the middleware for basic authentication
func BasicAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		// Parse basic auth header
		if !strings.HasPrefix(auth, "Basic ") {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		payload, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		username := pair[0]
		// Store username in context for later use
		c.Locals("username", username)
		return c.Next()
	}
}

// BasicAuthMiddlewareAdmin is the middleware for admin basic authentication for admin routes
func BasicAuthMiddlewareAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		// Parse basic auth header
		if !strings.HasPrefix(auth, "Basic ") {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		payload, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		username, password := pair[0], pair[1]
		if checkCredentials(username, password) {
			// Store username in context for later use
			c.Locals("username", username)
			return c.Next()
		} else {
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}
	}
}

func checkCredentials(username, password string) bool {
	credentials := authorized_credentials
	storedPass, ok := credentials[username]
	if !ok {
		return false
	}
	// Compare the provided password with the stored password
	return password == storedPass
}
