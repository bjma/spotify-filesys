package filesys

import (
	"fmt"

	// Libraries
	"github.com/zmb3/spotify"
)

var client *spotify.Client = nil

// Can be a playlist or folder
// If format is folder, children array is populated with SimplePlaylist pointers
// If format is playlist, children array is populated with PlaylistTrack pointers
type Node struct {
	Name 	 	 string 	   // Playlist or folder name
	Format 	 	 string   	   // {"folder", "playlist", "track"}
	Children 	 []Node 	   // Children of node
	Num_children int
}

// Should be able to support the following commands:
// - tree
// - cwp
// - mv
type Tree struct {
	client 		 *spotify.Client 		 // Reference to Spotify client for mounting purposes
	cwp 	     *spotify.SimplePlaylist // Current working playlist
	children 	 []Node  				 // Subdirectories
	num_children int				
}

func contains(list []Node, comparator string) bool {
	for _, elem := range list {
		if elem.Name == comparator {
			return true
		}
	}
	return false
}

// Constructs a folder by doing an exhaustive search on user's playlists
// If folder name matches parameter, then add to children and return array
// This is really fucking slow lmao
func constructFolder(dirname string, index int, folders map[string]string) Node {
	offset := index
	limit := 20

	var children []Node
	flag := false

	for {
		playlists, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{ Offset: &offset, Limit: &limit })
		if err != nil || playlists == nil || flag {
			break
		}
		for _, playlist := range playlists.Playlists {
			playlist_uri := string(playlist.URI)

			// Construct children
			if folders[playlist_uri] == dirname {
				// Append Node of format "Playlist"
				children = append(children, constructPlaylist(playlist.Name, playlist.ID))
			} else {
				flag = true
			}
		}
	}
	return Node{
		Name: dirname,
		Format: "folder",
		Children: children,
		Num_children: len(children),
	}
}

// Constructs a Node of format "Playlist"
func constructPlaylist(name string, playlist_id spotify.ID) Node {
	//var children []Node

	/*tracks, err := client.GetPlaylistTracks(playlist_id)
	if err != nil {
		panic(err)
	}

	for _, track := range tracks.Tracks {
		children = append(children, Node{
			Name: track.Track.Name,
			Format: "track",
			Children: nil,
			Num_children: 0,
		})
	}*/

	return Node{
		Name: name,
		Format: "playlist",
		Children: nil,
		Num_children: 0,
	}
}

// Parses entire user library for building filesystem
func parseLibrary(folders map[string]string) ([]Node, int) {	
	var nodes []Node

	offset := 0
	limit := 20
	count := 0

	// Iterates through user's entire playlist library and initializes data
	// needed to generate directory tree
	for {
		playlists, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{ Offset: &offset, Limit: &limit })
		if err != nil || playlists == nil || len(playlists.Playlists) < 1 {
			break
		}

		for _, playlist := range playlists.Playlists {
			playlist_uri := string(playlist.URI)

			// If playlist belongs to folder, construct folder and set it to children;
			// Else, simply append it.
			if folders[playlist_uri] == "" { 
				nodes = append(nodes, constructPlaylist(playlist.Name, playlist.ID))
				count++
			} else { 						
				node := constructFolder(folders[playlist_uri], count, folders)
				if !contains(nodes, node.Name) {
					nodes = append(nodes, constructFolder(folders[playlist_uri], count, folders))
				}
			}
		}
		// Increment offset for pagination
		offset += limit
	}
	return nodes, count
}

// Debugging purposes
func PrintTree(t *Tree) {
	tree := t.children

	fmt.Println(".")
	for _, node := range tree {
		fmt.Printf("\t%s\n", node.Name)
		if (node.Format == "folder") {
			for _, child := range node.Children {
				fmt.Print("\t\t")
				fmt.Println(child.Name)
			}
		}
	}
}

// Builds a directory tree from client
// We might need to parse folders somehow
func BuildTree(c *spotify.Client, f map[string]string) *Tree {
	client = c

	nodes, _ := parseLibrary(f)
	
	return &Tree{
		client: client,
		cwp: nil, // figure out what to do for root
		children: nodes,
	}
}
