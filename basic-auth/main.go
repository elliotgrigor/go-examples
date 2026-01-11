package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// NOTE: Colon `:` must be an illegal character in usernames and passwords

const (
	authHeaderPrefix = "Basic "
	cookieName       = "basicauth_session"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /secret", handleSecret)
	mux.HandleFunc("POST /auth/login", handleLogin)
	mux.HandleFunc("GET /auth/logout", handleLogout)

	http.ListenAndServe(":3000", mux)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func handleSecret(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("secret.html"))
	tmpl.Execute(w, nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	authz := r.Header.Get("Authorization")
	if authz == "" || !strings.HasPrefix(authz, authHeaderPrefix) {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	b64auth := strings.TrimPrefix(authz, authHeaderPrefix)
	b, err := base64.StdEncoding.DecodeString(b64auth)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	creds := strings.SplitN(string(b), ":", 2)
	if !verifyCredentials(creds) {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Println("verifyCredentials: credentials contains an illegal character")
		return
	}

	// TODO: Validate user and password

	sessionId := rand.Text()
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true when HTTPS-enabled
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
}

func verifyCredentials(credentials []string) bool {
	for _, c := range credentials {
		if strings.Contains(c, ":") {
			return false
		}
	}
	return true
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	_ = c
	// TODO: Find session and delete from session store
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
