package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const accessTokenKey contextKey = "access_token"

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func CORSMiddleware(frontendURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", frontendURL)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				slog.Debug("CORS preflight request", "path", r.URL.Path, "origin", r.Header.Get("Origin"))
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		slog.Info("HTTP request started", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery, "remote_addr", r.RemoteAddr)

		next.ServeHTTP(rr, r)

		slog.Info("HTTP request completed", "method", r.Method, "path", r.URL.Path, "status", rr.statusCode, "duration_ms", time.Since(start).Milliseconds())
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			slog.Warn("Unauthorized request - missing or empty access_token cookie", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		slog.Debug("Request authenticated via cookie", "path", r.URL.Path)

		ctx := context.WithValue(r.Context(), accessTokenKey, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccessTokenFromContext(ctx context.Context) string {
	token, _ := ctx.Value(accessTokenKey).(string)
	return token
}
