package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const errInternalServerError = "Internal Server Error"

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
		// TODO: A lot of this can be done in an [Authenticator] receiver function

		// Exchange code for a token
		token, err := auth.Exchange(r.Context(), r.URL.Query().Get("code")) // Use state here?
		if err != nil {
			http.Error(w, errInternalServerError, http.StatusInternalServerError)
			log.Println("handleCallback: failed to exchange code for token")
		}

		rawIdToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, errInternalServerError, http.StatusInternalServerError)
			log.Println("handleCallback: failed to extract raw id token")
		}

		oidcCfg := &oidc.Config{
			ClientID: auth.ClientID,
		}
		idToken, err := auth.Verifier(oidcCfg).Verify(r.Context(), rawIdToken)
		if err != nil {
			http.Error(w, errInternalServerError, http.StatusInternalServerError)
			log.Println("handleCallback: failed to verify id token")
		}

		var profile map[string]any
		idToken.Claims(&profile)

		fmt.Println(profile)
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
