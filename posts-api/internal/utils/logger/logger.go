package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger
type Logger struct {
	Logger *zap.Logger
}

// NewLogger creates a new logger
func NewLogger(level, format, output, file string) (*Logger, error) {
	// Configure logger
	var config zap.Config
	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		config.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Set output
	if output == "file" && file != "" {
		config.OutputPaths = []string{file}
		config.ErrorOutputPaths = []string{file}

		// Create directory if it doesn't exist
		dir := file[:len(file)-len("/post-api.log")]
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}
		}
	} else {
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}
	}

	// Build logger
	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.Logger.Sugar().Debugw(msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.Logger.Sugar().Infow(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.Logger.Sugar().Warnw(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	if err != nil {
		fields = append(fields, "error", err.Error())
	}
	l.Logger.Sugar().Errorw(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...interface{}) {
	if err != nil {
		fields = append(fields, "error", err.Error())
	}
	l.Logger.Sugar().Fatalw(msg, fields...)
}

// Panic logs a panic message and panics
func (l *Logger) Panic(msg string, err error, fields ...interface{}) {
	if err != nil {
		fields = append(fields, "error", err.Error())
	}
	l.Logger.Sugar().Panicw(msg, fields...)
}