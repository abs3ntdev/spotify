package spotify

import (
	"context"
	"fmt"
	"time"
)

type TrackExternalIDs struct {
	ISRC string `json:"isrc"`
	EAN  string `json:"ean"`
	UPC  string `json:"upc"`
}

// SimpleTrack contains basic info about a track.
type SimpleTrack struct {
	Album   SimpleAlbum    `json:"album"`
	Artists []SimpleArtist `json:"artists"`
	// The disc number (usually 1 unless the album consists of more than one disc).
	DiscNumber int `json:"disc_number"`
	// The length of the track, in milliseconds.
	Duration int `json:"duration_ms"`
	// Whether or not the track has explicit lyrics.
	// true => yes, it does; false => no, it does not.
	Explicit bool `json:"explicit"`
	// External URLs for this track.
	ExternalURLs map[string]string `json:"external_urls"`
	// ExternalIDs are IDs for this track in other databases
	ExternalIDs TrackExternalIDs `json:"external_ids"`
	// A link to the Web API endpoint providing full details for this track.
	Endpoint string `json:"href"`
	ID       ID     `json:"id"`
	Name     string `json:"name"`
	// A URL to a 30 second preview (MP3) of the track.
	PreviewURL string `json:"preview_url"`
	// The number of the track.  If an album has several
	// discs, the track number is the number on the specified
	// DiscNumber.
	TrackNumber int `json:"track_number"`
	URI         URI `json:"uri"`
	// Type of the track
	Type string `json:"type"`
}

func (st SimpleTrack) String() string {
	return fmt.Sprintf("TRACK<[%s] [%s]>", st.ID, st.Name)
}

// FullTrack provides extra track data in addition to what is provided by SimpleTrack.
type FullTrack struct {
	SimpleTrack
	// The album on which the track appears. The album object includes a link in href to full information about the album.
	Album SimpleAlbum `json:"album"`
	// Known external IDs for the track.
	ExternalIDs map[string]string `json:"external_ids"`

	// IsPlayable defines if the track is playable. It's reported when the "market" parameter is passed to the tracks
	// listing API.
	// See: https://developer.spotify.com/documentation/general/guides/track-relinking-guide/
	IsPlayable *bool `json:"is_playable"`
}

// PlaylistTrack contains info about a track in a playlist.
type PlaylistTrack struct {
	// The date and time the track was added to the playlist.
	// You can use the TimestampLayout constant to convert
	// this field to a time.Time value.
	// Warning: very old playlists may not populate this value.
	AddedAt string `json:"added_at"`
	// The Spotify user who added the track to the playlist.
	// Warning: vary old playlists may not populate this value.
	AddedBy User `json:"added_by"`
	// Whether this track is a local file or not.
	IsLocal bool `json:"is_local"`
	// Information about the track.
	Track FullTrack `json:"track"`
}

// SavedTrack provides info about a track saved to a user's account.
type SavedTrack struct {
	// The date and time the track was saved, represented as an ISO
	// 8601 UTC timestamp with a zero offset (YYYY-MM-DDTHH:MM:SSZ).
	// You can use the TimestampLayout constant to convert this to
	// a time.Time value.
	AddedAt   string `json:"added_at"`
	FullTrack `json:"track"`
}

// TimeDuration returns the track's duration as a time.Duration value.
func (t *SimpleTrack) TimeDuration() time.Duration {
	return time.Duration(t.Duration) * time.Millisecond
}

// GetTrack gets Spotify catalog information for
// a single track identified by its unique Spotify ID.
//
// API Doc: https://developer.spotify.com/documentation/web-api/reference/tracks/get-track/
//
// Supported options: Market
func (c *Client) GetTrack(ctx context.Context, id ID, opts ...RequestOption) (*FullTrack, error) {
	spotifyURL := c.baseURL + "tracks/" + string(id)

	var t FullTrack

	if params := processOptions(opts...).urlParams.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	err := c.get(ctx, spotifyURL, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
