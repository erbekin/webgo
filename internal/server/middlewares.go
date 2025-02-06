package server

import (
	"log"
	"net/http"
	"time"
)

// Middleware: log requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("INFO: %s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("INFO: İşlem süresi: %v", duration)
	})
}
