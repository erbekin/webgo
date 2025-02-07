package server

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	username = "admin"
	password = "admin"
)

// Middleware: log requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("INFO: %s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("INFO: İşlem süresi: %v", duration)
		fmt.Println("---------------------------------------")
	})
}

func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get path
		path := r.URL.Path
		if strings.HasPrefix(path, "/auth/") || path == "/auth" {

			isAuthorized := authUser(r)
			if isAuthorized {
				usr := BasicUser{
					Username: username,
					Password: password,
					Auth:     true,
				}
				ctx := context.WithValue(r.Context(), "auth", usr)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authUser(r *http.Request) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}

	usernameHash := sha256.Sum256([]byte(u))
	passwordHash := sha256.Sum256([]byte(p))
	expectedUsername := sha256.Sum256([]byte(username))
	expectedPassword := sha256.Sum256([]byte(password))
	usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsername[:]) == 1
	passwordMatch := subtle.ConstantTimeCompare(expectedPassword[:], passwordHash[:]) == 1
	if usernameMatch && passwordMatch {
		return true

	} else {
		return false
	}
}
