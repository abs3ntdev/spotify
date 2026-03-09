// Command profile gets the profile information about the current Spotify user.
package main

import (
	"context"
	"fmt"
	spotifyauth "github.com/abs3ntdev/spotify/v2/auth"
	"log"
	"os"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/abs3ntdev/spotify/v2"
)

func main() {
	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	user, err := client.CurrentUser(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("Display name:", user.DisplayName)
	fmt.Println("Spotify URI:", string(user.URI))
	fmt.Println("Endpoint:", user.Endpoint)
}
