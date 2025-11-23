package middleware

import (
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// SentryMiddleware создает middleware для интеграции с Sentry
func SentryMiddleware() func(http.Handler) http.Handler {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		// Если DSN не указан, возвращаем noop middleware
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      os.Getenv("ENV"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		// Если инициализация не удалась, возвращаем noop
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	handler := sentryhttp.New(sentryhttp.Options{
		Repanic:         true,
		WaitForDelivery: false,
	})

	return handler.Handle
}
