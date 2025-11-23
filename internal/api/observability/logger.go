package observability

import (
	"context"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
	With(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

type zapLogger struct {
	logger *zap.Logger
}

func NewLogger() (Logger, error) {
	var config zap.Config
	env := os.Getenv("ENV")
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

func (l *zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	zapFields := l.convertFields(fields)
	l.logger.Info(msg, zapFields...)
}

func (l *zapLogger) Error(ctx context.Context, msg string, err error, fields ...Field) {
	zapFields := l.convertFields(fields)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	l.logger.Error(msg, zapFields...)
}

func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	zapFields := l.convertFields(fields)
	l.logger.Warn(msg, zapFields...)
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	zapFields := l.convertFields(fields)
	l.logger.Debug(msg, zapFields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	zapFields := l.convertFields(fields)
	return &zapLogger{logger: l.logger.With(zapFields...)}
}

func (l *zapLogger) convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}

// Fallback logger для случаев когда zap недоступен
type fallbackLogger struct{}

func NewFallbackLogger() Logger {
	return &fallbackLogger{}
}

func (l *fallbackLogger) Info(ctx context.Context, msg string, fields ...Field) {
	log.Printf("[INFO] %s", msg)
}

func (l *fallbackLogger) Error(ctx context.Context, msg string, err error, fields ...Field) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", msg, err)
	} else {
		log.Printf("[ERROR] %s", msg)
	}
}

func (l *fallbackLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	log.Printf("[WARN] %s", msg)
}

func (l *fallbackLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	log.Printf("[DEBUG] %s", msg)
}

func (l *fallbackLogger) With(fields ...Field) Logger {
	return l
}

