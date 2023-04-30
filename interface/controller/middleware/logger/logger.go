package logger

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

func WithLogger(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, loggerKey{}, logrus.New())
	return ctx
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), loggerKey{}, logrus.New()))
		next.ServeHTTP(w, r)
	})
}

func FromContext(ctx context.Context) *logrus.Logger {
	log, ok := ctx.Value(loggerKey{}).(*logrus.Logger)
	if !ok {
		return nil
	}
	return log
}
