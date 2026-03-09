package spotify

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestLibraryContains(t *testing.T) {
	client, server := testClientString(http.StatusOK, `[ false, true ]`)
	defer server.Close()

	contains, err := client.LibraryContains(context.Background(), "spotify:track:0udZHhCi7p1YzMlvI4fXoK", "spotify:track:55nlbqqFVnSsArIeYSQlqx")
	if err != nil {
		t.Error(err)
	}
	if l := len(contains); l != 2 {
		t.Error("Expected 2 results, got", l)
	}
	if contains[0] || !contains[1] {
		t.Error("Expected [false, true], got", contains)
	}
}

func TestSaveToLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.SaveToLibrary(context.Background(), "spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}

func TestSaveToLibraryFailure(t *testing.T) {
	client, server := testClientString(http.StatusUnauthorized, `
{
  "error": {
    "status": 401,
    "message": "Invalid access token"
  }
}`)
	defer server.Close()
	err := client.SaveToLibrary(context.Background(), "spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M")
	if err == nil {
		t.Error("Expected error and didn't get one")
	}
}

func TestSaveToLibraryWithContextCancelled(t *testing.T) {
	client, server := testClientString(http.StatusOK, ``)
	defer server.Close()

	ctx, done := context.WithCancel(context.Background())
	done()

	err := client.SaveToLibrary(ctx, "spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M")
	if !errors.Is(err, context.Canceled) {
		t.Error("Expected error and didn't get one")
	}
}

func TestRemoveFromLibrary(t *testing.T) {
	client, server := testClientString(http.StatusOK, "")
	defer server.Close()

	err := client.RemoveFromLibrary(context.Background(), "spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M")
	if err != nil {
		t.Error(err)
	}
}
