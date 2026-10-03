package log

import (
	"context"
	stdlog "log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Client interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
	With(fields ...Field) Client
}

type zapClient struct {
	logger *zap.Logger
}

func NewClient() (Client, error) {
	var config zap.Config
	if os.Getenv("ENV") == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &zapClient{logger: logger}, nil
}

func NewFallbackClient() Client {
	return &fallbackClient{}
}

func (c *zapClient) Info(ctx context.Context, msg string, fields ...Field) {
	c.logger.Info(msg, c.convertFields(fields)...)
}

func (c *zapClient) Error(ctx context.Context, msg string, err error, fields ...Field) {
	zapFields := c.convertFields(fields)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	c.logger.Error(msg, zapFields...)
}

func (c *zapClient) Warn(ctx context.Context, msg string, fields ...Field) {
	c.logger.Warn(msg, c.convertFields(fields)...)
}

func (c *zapClient) Debug(ctx context.Context, msg string, fields ...Field) {
	c.logger.Debug(msg, c.convertFields(fields)...)
}

func (c *zapClient) With(fields ...Field) Client {
	return &zapClient{logger: c.logger.With(c.convertFields(fields)...)}
}

func (c *zapClient) convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}

type fallbackClient struct{}

func (c *fallbackClient) Info(ctx context.Context, msg string, fields ...Field) {
	stdlog.Printf("[INFO] %s", msg)
}

func (c *fallbackClient) Error(ctx context.Context, msg string, err error, fields ...Field) {
	if err != nil {
		stdlog.Printf("[ERROR] %s: %v", msg, err)
		return
	}
	stdlog.Printf("[ERROR] %s", msg)
}

func (c *fallbackClient) Warn(ctx context.Context, msg string, fields ...Field) {
	stdlog.Printf("[WARN] %s", msg)
}

func (c *fallbackClient) Debug(ctx context.Context, msg string, fields ...Field) {
	stdlog.Printf("[DEBUG] %s", msg)
}

func (c *fallbackClient) With(fields ...Field) Client {
	return c
}
