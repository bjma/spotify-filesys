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

func contains(list []string, comparator string) bool {
	for _, elem := range list {
		if elem == comparator {
			return true
		}
	}
	return false
}

// Fetches playlists within entire user library
func FetchPlaylists() []spotify.SimplePlaylist {
	var ret []spotify.SimplePlaylist
	offset := 0
	limit := 20

	for {
		playlists, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{ Offset: &offset, Limit: &limit })
		if err != nil || playlists == nil || len(playlists.Playlists) < 1 {
			break
		}

		for _, playlist := range playlists.Playlists {
			ret = append(ret, playlist)
		}
		// Increment offset for pagination
		offset += limit
	}
	return ret
}

// Constructs a folder by doing an exhaustive search on user's playlists
// If folder name matches parameter, then add to children and return array
// This is really fucking slow lmao
func constructFolder(playlists []spotify.SimplePlaylist, dirname string, index int, folders map[string]string) (Node, int) {
	var children []Node
	iter := 0

	for i := index; i < len(playlists) && folders[string(playlists[i].URI)] == dirname; i++ {
		playlist := playlists[i]
		children = append(children, constructPlaylist(playlist.Name, playlist.ID))
		iter++
	}
	node := Node{
		Name: dirname,
		Format: "folder",
		Children: children,
		Num_children: len(children),
	}
	return node, iter
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
func parseLibrary(folders map[string]string) []Node {	
	var nodes []Node


	playlists := FetchPlaylists()
	//fmt.Printf("Expected 59 got %d\n", len(playlists))

	i := 0
	// Iterates through user's entire playlist library and initializes data needed to generate directory tree
	for i < len(playlists) {
		uri := string(playlists[i].URI)
		// If current playlist is in a folder, parse it as a folder, then append to playlist
		// Else, just append to nodes as playlist
		if folders[uri] != "" {
			node, iter := constructFolder(playlists, folders[uri], i, folders)
			nodes = append(nodes, node)
			i += iter
		} else {
			nodes = append(nodes, constructPlaylist(playlists[i].Name, playlists[i].ID))
			i++
		}
	}
	
	return nodes
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

	nodes := parseLibrary(f)
	
	return &Tree{
		client: client,
		cwp: nil, // figure out what to do for root
		children: nodes,
	}
}
