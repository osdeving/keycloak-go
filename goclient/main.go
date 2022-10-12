package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)



const (
	clientID = "myclient"
	clientSecret = "QIc50AM1Spw9qE4tXBMtoN5wr0qnXkoh"
)

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/realms/myrealm")

	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	// o certo é deixar dinâmico
	state := "123"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "State Inválido", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, r.URL.Query().Get("code"))

		if err != nil {
			http.Error(w, "Falha ao trocar o token", http.StatusInternalServerError)
			return
		}


		idToken, ok := token.Extra("id_token").(string)

		if !ok {
			http.Error(w, "Falha ao gerar IDToken", http.StatusInternalServerError)
			return
		}

		resp := struct {
			AccessToken *oauth2.Token
			IDToken string
		}{
			token,
			idToken,
		}

		data, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}