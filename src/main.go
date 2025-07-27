package main

import (
	"pg_bridge_go/config"
	"pg_bridge_go/database"
	"pg_bridge_go/logger"
	"pg_bridge_go/routes"
)

func init() {
	// Initialize logger
	logger.Init(true)
	defer logger.Close()

	// Initialize environment configuration
	logger.Info("Initializing environment configuration")
	config.InitEnvConfig()

	// initialize SetupDatabase
	logger.Info("Setting up database")
	database.SetupDatabase()

	// load default callback url
	logger.Info("Load default callback url")
	config.LoadDefaultCallbackUrl()

	// Load the encryption key
	logger.Info("Loading encryption key")
	config.LoadEncryptionKey()

	// Load App Port
	logger.Info("Loading application port")
	config.LoadAppPort()
}

// Entrypoint for app fiber.
func main() {
	// Load the routes
	r := routes.SetupRouter()

	// Start the HTTP API
	r.Listen("0.0.0.0:" + config.AppPort)
}
