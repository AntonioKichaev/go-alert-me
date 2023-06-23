package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lw := loggingResponseWriter{
			ResponseWriter: w,
			Rd:             &responseData{Status: 0, Size: 0},
		}
		t := time.Now()
		next.ServeHTTP(&lw, r)
		c := time.Since(t).String()
		Log.Info("HTTP request",
			zap.String("path", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", lw.Rd.Status),
			zap.Int("length", lw.Rd.Size),
			zap.String("duration", c),
		)

	})
}
