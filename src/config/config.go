package config

import (
	"encoding/hex"
	"log"
	"os"
	"pg_bridge_go/global_var"
	"pg_bridge_go/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	MasterKey   []byte
	CallbackUrl string
	AppPort     string
)

// InitEnvConfig loads environment variables from .env file
func InitEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("failed to load .env file", zap.Error(err))
		log.Panic("Failed to load .env file:", err)
	}
}

func GetEnvDatabase() global_var.DatabaseConnection {
	return global_var.DatabaseConnection{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
	}
}

// LoadEncryptionKey parses and stores the master key from env
func LoadEncryptionKey() {
	keyHex := os.Getenv("MASTER_KEY")
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		logger.Error("invalid MASTER_KEY hex in env", zap.Error(err))
		log.Panic("Invalid MASTER_KEY hex:", err)
	}
	if len(key) != 32 {
		log.Panicf("Master key must be 32 bytes (64 hex chars) for AES-256, got %d bytes", len(key))
	}
	MasterKey = key
}

func LoadDefaultCallbackUrl() {
	CallbackUrl = os.Getenv("DEFAULT_CALLBACK")
}

func LoadAppPort() {
	AppPort = os.Getenv("APP_PORT")
}

func GetEncryptionKey() []byte {
	return MasterKey
}
