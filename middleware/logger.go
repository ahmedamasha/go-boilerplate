package middleware

import (
    "log"
    "net/http"
    "time"
)

// Logger is a middleware that logs the incoming HTTP request and response status.
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)

        // Wrap the ResponseWriter to capture status code
        lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(lrw, r)

        log.Printf(
            "Completed %s %s with %d in %v",
            r.Method,
            r.URL.Path,
            lrw.statusCode,
            time.Since(start),
        )
    })
}

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}