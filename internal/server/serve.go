package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Serve starts server. Server will listen on addr
func Serve(addr string) error {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if _, err := fmt.Fprint(w, "<h1> Welcome </h1>"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		response := response{
			Message: "Nothing to see here.",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	// Use middleware
	loggedMux := loggingMiddleware(mux)

	// Start
	log.Println("[Listening]" + addr)

	return http.ListenAndServe(addr, loggedMux)
}

type response struct {
	Message string `json:"message"`
}
