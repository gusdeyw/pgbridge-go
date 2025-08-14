package middleware

import (
	"crypto/subtle"
	"os"
	"pg_bridge_go/helper"
)

// authorized_credentials stores username->hashed_password pairs
// Passwords should be bcrypt hashed, not plain text
// In production, consider moving these to environment variables or a secure database
var authorized_credentials = map[string]string{
	// Default admin user with hashed password (change in production)
	// Default password: "admin123" - CHANGE THIS IN PRODUCTION!
	"admin": "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // bcrypt hash of "admin123"
}

// initializeCredentialsFromEnv loads credentials from environment variables if available
func initializeCredentialsFromEnv() {
	// Load admin credentials from environment if available
	if adminUser := os.Getenv("ADMIN_USERNAME"); adminUser != "" {
		if adminPassHash := os.Getenv("ADMIN_PASSWORD_HASH"); adminPassHash != "" {
			authorized_credentials[adminUser] = adminPassHash
		}
	}
}

// secureComparePasswords performs constant-time comparison to prevent timing attacks
func secureComparePasswords(provided, stored string) bool {
	// Use bcrypt to verify hashed passwords
	return helper.VerifyPassword(provided, stored)
}

// constantTimeStringCompare provides constant-time string comparison for additional security
func constantTimeStringCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
