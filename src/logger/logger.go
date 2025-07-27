package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log   *zap.Logger        // Structured logger
	Sugar *zap.SugaredLogger // Sugared logger
)

// Init initializes zap logger.
// If debug == true, uses Development config.
// If debug == false, uses Production config.
func Init(debug bool) {
	var err error

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create a file writer
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	writeSyncer := zapcore.AddSync(file)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		writeSyncer,
		zap.InfoLevel, // Minimum log level
	)

	Log = zap.New(core, zap.AddCaller())
	Sugar = Log.Sugar()
	Log.Info("Logger initialized with file output")
}

func Close() {
	_ = Log.Sync()
}

// === READY-TO-USE SIMPLE FUNCTIONS ===

// Info logs an info message.
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Debug logs a debug message.
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Warn logs a warning.
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error logs an error.
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// SugarInfo logs using SugaredLogger (printf style).
func SugarInfo(template string, args ...interface{}) {
	Sugar.Infof(template, args...)
}

// SugarWarn logs using SugaredLogger (printf style).
func SugarWarn(template string, args ...interface{}) {
	Sugar.Warnf(template, args...)
}

// SugarError logs using SugaredLogger (printf style).
func SugarError(template string, args ...interface{}) {
	Sugar.Errorf(template, args...)
}
