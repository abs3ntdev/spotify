package main

import (
	"context"
	"flag"
	"log"

	"golang.org/x/oauth2"

	"github.com/abs3ntdev/spotify/v2"
	spotifyauth "github.com/abs3ntdev/spotify/v2/auth"
)

var auth = spotifyauth.New(spotifyauth.WithRedirectURL("http://localhost:3000/login_check"))

func main() {
	code := flag.String("code", "", "authorization code to negotiate by token")
	flag.Parse()

	if *code == "" {
		log.Fatal("code required")
	}
	if err := authorize(*code); err != nil {
		log.Fatal("error while negotiating the token: ", err)
	}
}

func authorize(code string) error {
	ctx := context.Background()
	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return err
	}
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	client := spotify.New(httpClient)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		return err
	}
	log.Printf("Logged in as %s\n", user.DisplayName)

	return nil
}
