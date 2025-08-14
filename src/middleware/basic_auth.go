package middleware

import (
	"encoding/base64"
	"pg_bridge_go/helper"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// rateLimiter provides basic rate limiting for authentication attempts
type rateLimiter struct {
	attempts map[string][]time.Time
	mutex    sync.RWMutex
}

var authRateLimiter = &rateLimiter{
	attempts: make(map[string][]time.Time),
}

// isRateLimited checks if the IP has exceeded the rate limit
func (rl *rateLimiter) isRateLimited(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-5 * time.Minute) // 5-minute window

	// Clean old attempts
	if attempts, exists := rl.attempts[ip]; exists {
		validAttempts := []time.Time{}
		for _, attempt := range attempts {
			if attempt.After(windowStart) {
				validAttempts = append(validAttempts, attempt)
			}
		}
		rl.attempts[ip] = validAttempts

		// Check if rate limited (max 5 attempts per 5 minutes)
		if len(validAttempts) >= 5 {
			return true
		}
	}

	return false
}

// recordAttempt records a failed authentication attempt
func (rl *rateLimiter) recordAttempt(ip string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	rl.attempts[ip] = append(rl.attempts[ip], now)
}

func init() {
	// Initialize credentials from environment on package load
	initializeCredentialsFromEnv()
}

// BasicAuthMiddleware is the middleware for basic authentication
// This version only validates the format but doesn't check credentials
// Use BasicAuthMiddlewareAdmin for actual credential validation
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
		// Validate username format (basic security check)
		if len(strings.TrimSpace(username)) == 0 {
			return helper.SendResponse(fiber.StatusUnauthorized, "Invalid credentials", nil, c)
		}

		// Store username in context for later use
		c.Locals("username", username)
		return c.Next()
	}
}

// BasicAuthMiddlewareAdmin is the middleware for admin basic authentication for admin routes
func BasicAuthMiddlewareAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		// Check rate limiting
		if authRateLimiter.isRateLimited(clientIP) {
			return helper.SendResponse(fiber.StatusTooManyRequests, "Too many authentication attempts. Please try again later.", nil, c)
		}

		auth := c.Get("Authorization")
		if auth == "" {
			authRateLimiter.recordAttempt(clientIP)
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		// Parse basic auth header
		if !strings.HasPrefix(auth, "Basic ") {
			authRateLimiter.recordAttempt(clientIP)
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		payload, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			authRateLimiter.recordAttempt(clientIP)
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			authRateLimiter.recordAttempt(clientIP)
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}

		username, password := pair[0], pair[1]
		
		// Validate input lengths to prevent potential attacks
		if len(username) > 255 || len(password) > 255 {
			authRateLimiter.recordAttempt(clientIP)
			return helper.SendResponse(fiber.StatusUnauthorized, "Invalid credentials", nil, c)
		}

		if checkCredentials(username, password) {
			// Store username in context for later use
			c.Locals("username", username)
			return c.Next()
		} else {
			authRateLimiter.recordAttempt(clientIP)
			return helper.SendResponse(fiber.StatusUnauthorized, "Not Authorized", nil, c)
		}
	}
}

func checkCredentials(username, password string) bool {
	credentials := authorized_credentials
	storedPassHash, ok := credentials[username]
	if !ok {
		// Always perform a dummy bcrypt operation to prevent timing attacks
		// even when the username doesn't exist
		helper.VerifyPassword("dummy", "$2a$10$dummy.hash.to.prevent.timing.attacks.dummy.hash.value")
		return false
	}
	
	// Use secure password comparison with bcrypt
	return secureComparePasswords(password, storedPassHash)
}
