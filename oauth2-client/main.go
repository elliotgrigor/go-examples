package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func NewAuthenticator() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/")
	if err != nil {
		return nil, err
	}

	cfg := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
		Endpoint:     provider.Endpoint(),
	}

	return &Authenticator{
		Provider: provider,
		Config:   cfg,
	}, nil
}

func main() {
	godotenv.Load()

	auth, err := NewAuthenticator()
	if err != nil {
		log.Fatalln(err)
	}

	state := rand.Text()
	nonce := rand.Text()
	stateB64 := base64.StdEncoding.EncodeToString([]byte(state))
	nonceB64 := base64.StdEncoding.EncodeToString([]byte(nonce))

	fmt.Println(auth.Config.AuthCodeURL(
		stateB64,
		oauth2.SetAuthURLParam("nonce", nonceB64)))
}
