# AGENTS.md — Coding Agent Instructions

## Project Overview

Go client library for the [Spotify Web API](https://developer.spotify.com/documentation/web-api).
Module path: `github.com/abs3ntdev/spotify/v2` (forked from `github.com/zmb3/spotify/v2`).
Minimum Go version: 1.26 (per `go.mod`), CI runs Go 1.26.

## Build / Lint / Test Commands

```bash
# Build all packages
go build ./...

# Run all tests (CI uses -race)
go test -race ./...

# Run a single test by name
go test -run TestGetAlbum -v ./...

# Run tests in a specific file's package
go test -v -run 'TestSaveToLibrary' .

# Lint (matches CI)
golangci-lint run ./...
```

## Project Structure

```
.                        # Root package "spotify" — all API types and client methods
├── auth/                # Sub-package "spotifyauth" — OAuth2 helpers and scope constants
├── examples/            # Runnable example programs (one per subdirectory)
├── test_data/           # JSON fixture files for tests (*.txt and *.json)
├── spotify.go           # Client, core types (URI, ID, Numeric, Error, Image, Followers)
├── request_options.go   # RequestOption functional options (Limit, Market, Offset, etc.)
├── page.go              # Paging types (basePage, FullArtistPage, SimpleAlbumPage, etc.)
├── album.go             # Album types + GetAlbum, GetAlbumTracks
├── artist.go            # Artist types + GetArtist, GetArtistAlbums
├── track.go             # Track types + GetTrack
├── playlist.go          # Playlist types + all playlist CRUD methods
├── user.go              # User types + CurrentUser, top items, followed artists, saved items
├── library.go           # Unified library endpoints: SaveToLibrary, RemoveFromLibrary, LibraryContains
├── show.go              # Show/Episode types + GetShow, GetShowEpisodes, GetEpisode
├── search.go            # Search method
├── player.go            # Player/playback control methods
└── ...                  # audio_analysis.go, audio_features.go, recommendation.go, etc.
```

## Code Style Guidelines

### Package & Imports

- The root package is `spotify`. The auth sub-package is `spotifyauth`.
- All source files use `package spotify` (tests use `package spotify` — not `_test`).
- Import groups are separated by blank lines in this order:
  1. Standard library
  2. Third-party (`golang.org/x/oauth2`, `github.com/stretchr/testify`)
  3. Internal (none currently — single-package library)
- Use `gofmt`/`goimports` formatting. No manual alignment beyond struct tags.

### Naming Conventions

- **Types**: PascalCase. Prefix with `Simple` or `Full` to distinguish API object variants
  (`SimpleAlbum`, `FullAlbum`, `SimpleTrack`, `FullTrack`).
- **Page types**: Named `{Type}Page` (e.g., `SimpleAlbumPage`, `SavedTrackPage`).
- **Client methods**: Named after the Spotify API operation, PascalCase, receiver `(c *Client)`.
  Examples: `GetAlbum`, `GetArtistAlbums`, `CreatePlaylist`, `SaveToLibrary`.
- **Functional options**: Top-level functions returning `RequestOption`
  (e.g., `Limit(n)`, `Market(code)`, `Offset(n)`).
- **Constants**: PascalCase for exported (`DateLayout`, `TimestampLayout`),
  camelCase for unexported (`defaultRetryDuration`).
- **Scopes**: Defined in `auth/auth.go` as `Scope*` constants (e.g., `ScopeUserLibraryModify`).

### Types & JSON

- **`ID`** (`string`): Base-62 Spotify identifier. Used in endpoint URL paths.
- **`URI`** (`string`): Full Spotify URI like `spotify:track:6rqhFgbbKwnb9MLmUQDhG6`.
- **`Numeric`** (`int`): Custom type with `UnmarshalJSON` that handles floats from the API.
  Use for any integer field that Spotify may return as a float (e.g., `disc_number`,
  `duration_ms`, `track_number`, `height`, `width`, follower `total`).
- JSON struct tags must match the Spotify API field names exactly (snake_case).
- Use `map[string]string` for `external_urls` and `external_ids` on full objects.
- Use `*bool` for nullable boolean fields (e.g., `IsExternallyHosted`).
- Removed API fields must be deleted from structs entirely — do not leave them
  commented out or with `omitempty`.

### Method Signatures

Every client method follows this pattern:

```go
func (c *Client) MethodName(ctx context.Context, requiredArgs, opts ...RequestOption) (ReturnType, error) {
    // 1. Build URL from c.baseURL + path
    // 2. Apply options via processOptions(opts...)
    // 3. Use c.get() for GET requests or c.execute() for PUT/POST/DELETE
    // 4. Return (result, nil) on success — NOT (result, err)
}
```

- Always take `context.Context` as the first parameter.
- Use variadic `...RequestOption` for optional query parameters.
- GET requests: use `c.get(ctx, url, &result)`.
- Mutating requests: build `*http.Request` with `http.NewRequestWithContext`, call `c.execute(req, result, statusCodes...)`.
- Return `nil` error on success, not the error variable from the successful decode.

### Error Handling

- Errors are prefixed with `"spotify: "` (e.g., `errors.New("spotify: this call supports 1 to 40 URIs")`).
- API errors decode into the `Error` type: `fmt.Sprintf("spotify: %s [%d]", e.Message, e.Status)`.
- Rate limit 429 responses populate `Error.RetryAfter` from the `Retry-After` header.
- Auto-retry logic is in `c.get()` and `c.execute()` — controlled by `WithRetry(true)`.

### Testing Patterns

- Tests use `testClientFile(statusCode, "test_data/filename.txt")` or
  `testClientString(statusCode, jsonString)` to create a mock HTTP client.
- Test data files live in `test_data/` as `.txt` or `.json` files containing raw JSON.
- Test functions follow `Test{MethodName}` naming (e.g., `TestGetAlbum`, `TestSaveToLibrary`).
- Tests use `testing.T` with `t.Error`/`t.Fatal` — the project also has `github.com/stretchr/testify`
  available but most tests use stdlib.
- Test fixtures should reflect current API responses (no removed fields).

### Spotify API Reference (Feb/March 2026)

This library tracks the **current** Spotify Web API as of the February 2026 and March 2026
changelogs. Key changes:

- **Removed endpoints**: `GET /albums`, `GET /artists`, `GET /tracks`, `GET /shows`,
  `GET /episodes`, `GET /audiobooks`, `GET /chapters`, `GET /artists/{id}/top-tracks`,
  `GET /artists/{id}/related-artists`, `GET /browse/new-releases`, `GET /browse/categories`,
  `GET /browse/categories/{id}`, `GET /users/{id}`, `GET /users/{id}/playlists`,
  `GET /markets`, all type-specific library save/remove/check endpoints,
  all follow/unfollow endpoints, old `/playlists/{id}/tracks` endpoints.
- **New unified library**: `PUT /me/library`, `DELETE /me/library`, `GET /me/library/contains`
  — all take `uris` query parameter (comma-separated Spotify URIs, max 40).
- **Playlist rename**: `tracks` -> `items` in all JSON responses and URL paths.
- **Search limit**: max 10 (was 50), default 5 (was 20).
- **Field removals**: `popularity`, `followers`, `available_markets`, `album_group`, `label`,
  `linked_from`, `country`, `email`, `product`, `explicit_content`, `publisher`.
- **March 2026 revert**: `external_ids` on Album and Track stays (was briefly removed).

When adding or modifying endpoints, always check the live API reference:
https://developer.spotify.com/documentation/web-api

### Do NOT

- Do not re-add removed fields or endpoints — the February 2026 changes are intentional.
- Do not use `interface{}` — use `any` or concrete types.
- Do not add dependencies without strong justification — this is a lightweight library.
- Do not put test helpers in non-test files.
- Do not use `println` or `log` in library code.
