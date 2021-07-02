package httpserver

import (
	"net/http"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			l := log.
				WithField("duration", time.Since(start).Milliseconds()).
				WithField("request", r.RequestURI).
				WithField("user_agent", r.UserAgent()).
				WithField("ip", r.RemoteAddr).
				WithField("proto", r.Proto)
			l.Info(r.Method)
		})
	}
}
