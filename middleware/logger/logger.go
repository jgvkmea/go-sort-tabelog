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

// TODO: ロガー生成する処理書く
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ここにミドルウェア処理を記載
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
