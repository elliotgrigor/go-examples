package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<head><link rel="icon" href="data:,"></head><h1>Hello</h1><a href="/auth/signin">Sign in</a>`)
}

func handleSignIn(auth *Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		state, nonce := rand.Text(), rand.Text()
		stateB64 := base64.StdEncoding.EncodeToString([]byte(state))
		nonceB64 := base64.StdEncoding.EncodeToString([]byte(nonce))
		url := auth.AuthCodeURL(stateB64,
			oauth2.SetAuthURLParam("nonce", nonceB64))

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func handleCallback(auth *Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = auth
		// TODO: Use authorization code to get access/id tokens
	}
}

func handleSignOut(w http.ResponseWriter, r *http.Request) {
	// TODO: Invalidate user session
}

func main() {
	godotenv.Load()

	auth, err := NewAuthenticator()
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /auth/signin", handleSignIn(auth))
	mux.HandleFunc("GET /auth/callback", handleCallback(auth))
	mux.HandleFunc("GET /auth/signout", handleSignOut)

	http.ListenAndServe(":3000", mux)
}
