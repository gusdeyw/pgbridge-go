package database

import (
	"fmt"
	"log"
	"os"
	"pg_bridge_go/config"
	"pg_bridge_go/db_var"
	"pg_bridge_go/global_var"
	"time"

	loggers "pg_bridge_go/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDatabase migrates and sets up the database.
func SetupDatabase() {
	var credentials global_var.DatabaseConnection = config.GetEnvDatabase()
	u := credentials.User
	p := credentials.Password
	n := credentials.DatabaseName
	// PostgreSQL DSN does not use charset/parseTime/loc like MySQL

	stdoutLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// PostgreSQL DSN format: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", credentials.Host, u, p, n, credentials.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: stdoutLogger,
	})
	if err != nil {
		loggers.Error("Could not connect to database", zap.Error(err))
		log.Panic("Could not connect to database:", err)
	}

	err = db.Debug().AutoMigrate(
		&db_var.UserT{},
		&db_var.PaymentGatewayCredentialT{},
		&db_var.PaymentGatewayTransactionT{},
	)
	if err != nil {
		loggers.Error("Error during migration", zap.Error(err))
		log.Panic("Error during migration:", err)
	}

	global_var.DB = db.Debug()
}
