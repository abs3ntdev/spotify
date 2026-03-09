package spotify

// Category is used by Spotify to tag items in.  For example, on the Spotify
// player's "Browse" tab.
type Category struct {
	// A link to the Web API endpoint returning full details of the category
	Endpoint string `json:"href"`
	// The category icon, in various sizes
	Icons []Image `json:"icons"`
	// The Spotify category ID.  This isn't a base-62 Spotify ID, its just
	// a short string that describes and identifies the category (ie "party").
	ID string `json:"id"`
	// The name of the category
	Name string `json:"name"`
}
