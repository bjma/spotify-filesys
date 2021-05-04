package filesys

import (
	_ "fmt"

	// Libraries
	"github.com/zmb3/spotify"
)

type Tree struct {
	client *spotify.Client 		// Reference to Spotify client for mounting purposes
	cwp *spotify.SimplePlaylist // Current working playlist
}