package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func getCredentialsFromHeader(r *http.Request) (credentials []string, errMsg string, status int) {
	authZ := r.Header.Get("Authorization")
	if authZ == "" || !strings.HasPrefix(authZ, authHeaderPrefix) {
		log.Println("[ERROR] getCredentialsFromHeader: Invalid authorization header format")
		return nil, "Invalid authorization header", http.StatusUnauthorized
	}

	b64AuthZ := strings.TrimPrefix(authZ, authHeaderPrefix)
	by, err := base64.StdEncoding.DecodeString(b64AuthZ)
	if err != nil {
		log.Println(err)
		return nil, "500 Internal Server Error", http.StatusInternalServerError
	}

	credentials = strings.SplitN(string(by), ":", 2)
	if !checkIllegalCharacters(credentials) {
		log.Println("[ERROR] getCredentialsFromHeader: Credential contains an illegal character")
		return nil, "500 Internal Server Error", http.StatusInternalServerError
	}

	return credentials, "", 0
}

func checkIllegalCharacters(credentials []string) bool {
	for _, c := range credentials {
		if strings.Contains(c, ":") {
			return false
		}
	}
	return true
}
