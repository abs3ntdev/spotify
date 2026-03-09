package spotify

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

func toURISlice(uris []URI) []string {
	result := make([]string, len(uris))
	for i, u := range uris {
		result[i] = string(u)
	}
	return result
}

// SaveToLibrary saves one or more items to the current user's library.
// Items are identified by their Spotify URIs (e.g. "spotify:track:4iV5W9uYEdYUVa79Axb7Rh").
// A maximum of 40 URIs can be sent per request.
// This call requires the ScopeUserLibraryModify scope.
func (c *Client) SaveToLibrary(ctx context.Context, uris ...URI) error {
	if l := len(uris); l == 0 || l > 40 {
		return errors.New("spotify: this call supports 1 to 40 URIs per call")
	}
	spotifyURL := c.baseURL + "me/library?uris=" + strings.Join(toURISlice(uris), ",")
	req, err := http.NewRequestWithContext(ctx, "PUT", spotifyURL, nil)
	if err != nil {
		return err
	}
	return c.execute(req, nil)
}

// RemoveFromLibrary removes one or more items from the current user's library.
// Items are identified by their Spotify URIs (e.g. "spotify:track:4iV5W9uYEdYUVa79Axb7Rh").
// A maximum of 40 URIs can be sent per request.
// This call requires the ScopeUserLibraryModify scope.
func (c *Client) RemoveFromLibrary(ctx context.Context, uris ...URI) error {
	if l := len(uris); l == 0 || l > 40 {
		return errors.New("spotify: this call supports 1 to 40 URIs per call")
	}
	spotifyURL := c.baseURL + "me/library?uris=" + strings.Join(toURISlice(uris), ",")
	req, err := http.NewRequestWithContext(ctx, "DELETE", spotifyURL, nil)
	if err != nil {
		return err
	}
	return c.execute(req, nil)
}

// LibraryContains checks if one or more items are saved in the current user's library.
// Items are identified by their Spotify URIs (e.g. "spotify:track:4iV5W9uYEdYUVa79Axb7Rh").
// A maximum of 40 URIs can be sent per request.
// The result is returned as a slice of bool values in the same order
// in which the URIs were specified.
func (c *Client) LibraryContains(ctx context.Context, uris ...URI) ([]bool, error) {
	if l := len(uris); l == 0 || l > 40 {
		return nil, errors.New("spotify: this call supports 1 to 40 URIs per call")
	}
	spotifyURL := c.baseURL + "me/library/contains?uris=" + strings.Join(toURISlice(uris), ",")

	var result []bool

	err := c.get(ctx, spotifyURL, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
