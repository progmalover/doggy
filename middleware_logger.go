package doggy

import (
	"context"
	"net/http"
	"time"

	"github.com/uber-go/zap"
	"github.com/urfave/negroni"
)

const zapKey = "zapKey"

func Logger(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	now := time.Now()
	log := LogFromContext(r.Context())
	ctx := ContextWithLog(r.Context(), log)
	ww := negroni.NewResponseWriter(rw)

	next(ww, r.WithContext(ctx))

	log.Info("Completed", zap.Float64("responsetime", time.Now().Sub(now).Seconds()),
		zap.String("path", r.URL.Path), zap.String("host", r.Host), zap.Int("code", ww.Status()))
}

func LogFromContext(ctx context.Context) zap.Logger {
	l, ok := ctx.Value(zapKey).(zap.Logger)
	if !ok {
		return zap.New(zap.NewJSONEncoder(TimeFormat("timestamp")), zap.AddCaller(), config.Logger.Level, zap.Output(config.Logger.File))
	}
	return l
}

func ContextWithLog(ctx context.Context, log zap.Logger) context.Context {
	return context.WithValue(ctx, zapKey, log)
}

func TimeFormat(key string) zap.TimeFormatter {
	return zap.TimeFormatter(func(t time.Time) zap.Field {
		return zap.String(key, t.Local().Format("2006-01-02T15:04:05.000Z0700"))
	})
}
