package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent clickjacking attacks
		c.Set("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// Enforce HTTPS (only set if not in development)
		// Note: Uncomment this for production with HTTPS
		// c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Control referrer information
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy - basic policy, adjust as needed
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")
		
		// Remove server information
		c.Set("Server", "")
		
		// Permissions Policy (formerly Feature Policy)
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		return c.Next()
	}
}