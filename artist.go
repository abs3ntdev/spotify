package spotify

import (
	"context"
	"fmt"
	"strings"
)

// SimpleArtist contains basic info about an artist.
type SimpleArtist struct {
	Name string `json:"name"`
	ID   ID     `json:"id"`
	// The Spotify URI for the artist.
	URI URI `json:"uri"`
	// A link to the Web API endpoint providing full details of the artist.
	Endpoint     string            `json:"href"`
	ExternalURLs map[string]string `json:"external_urls"`
}

// FullArtist provides extra artist data in addition to what is provided by [SimpleArtist].
type FullArtist struct {
	SimpleArtist
	// A list of genres the artist is associated with.  For example, "Prog Rock"
	// or "Post-Grunge".  If not yet classified, the slice is empty.
	Genres []string `json:"genres"`
	// Images of the artist in various sizes, widest first.
	Images []Image `json:"images"`
}

// GetArtist gets Spotify catalog information for a single artist, given its Spotify ID.
func (c *Client) GetArtist(ctx context.Context, id ID) (*FullArtist, error) {
	spotifyURL := fmt.Sprintf("%sartists/%s", c.baseURL, id)

	var a FullArtist
	err := c.get(ctx, spotifyURL, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetArtistAlbums gets Spotify catalog information about an [artist's albums].
//
// The AlbumType argument can be used to find particular types of album.
// If the Market is not specified, Spotify will likely return a lot
// of duplicates (one for each market in which the album is available).
//
// Supported options: [Market], [Limit], [Offset].
//
// [artist's albums]: https://developer.spotify.com/documentation/web-api/reference/get-an-artists-albums
func (c *Client) GetArtistAlbums(ctx context.Context, artistID ID, ts []AlbumType, opts ...RequestOption) (*SimpleAlbumPage, error) {
	spotifyURL := fmt.Sprintf("%sartists/%s/albums", c.baseURL, artistID)
	// add optional query string if options were specified
	values := processOptions(opts...).urlParams

	if ts != nil {
		types := make([]string, len(ts))
		for i := range ts {
			types[i] = ts[i].encode()
		}
		values.Set("include_groups", strings.Join(types, ","))
	}

	if query := values.Encode(); query != "" {
		spotifyURL += "?" + query
	}

	var p SimpleAlbumPage

	err := c.get(ctx, spotifyURL, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
