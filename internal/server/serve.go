package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

// Serve starts server. Server will listen on addr
func Serve(addr string) error {

	router := httprouter.New()
	// GET /
	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := fmt.Fprint(w, "Hello, World!"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// GET /echo
	router.HandlerFunc(http.MethodGet, "/echo", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		encoded := base64.StdEncoding.EncodeToString(body)
		cookiesJson, _ := json.Marshal(r.Cookies())
		resp := echoResponse{
			Method:  r.Method,
			Proto:   r.Proto,
			Headers: r.Header,
			Body:    encoded,
			Cookies: string(cookiesJson),
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	})

	// GET /auth
	router.HandlerFunc(http.MethodGet, "/auth", func(w http.ResponseWriter, r *http.Request) {
		ref := r.URL.Query().Get("ref")
		if ref == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, ref, http.StatusFound)
	})

	// GET /secret
	router.HandlerFunc(http.MethodGet, "/secret", func(w http.ResponseWriter, r *http.Request) {
		isAuthorized := authUser(r)

		if !isAuthorized {
			http.Redirect(w, r, "/auth?ref=/secret", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(w).Encode(struct {
			Message string    `json:"message"`
			User    BasicUser `json:"user-data"`
		}{
			Message: "Here is your secret",
			User: BasicUser{
				Username: username,
				Password: password,
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	})

	// Static
	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	// Use middleware
	loggedMux := loggingMiddleware(router)
	authMux := basicAuthMiddleware(loggedMux)

	// Start
	log.Println("[Listening]" + addr)

	return http.ListenAndServe(addr, authMux)
}

type response struct {
	Message string `json:"message"`
}

type echoResponse struct {
	Proto   string      `json:"proto"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
	Cookies string      `json:"cookies"`
}

type BasicUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     bool   `json:"auth"`
}
