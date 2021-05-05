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
	Name         string     // Playlist or folder name
	Format       string     // {"folder", "playlist", "track"}
	Children     []Node     // Children of node
	Id           spotify.ID // If playlist, we save its ID for access to tracks
	Num_children int
}

// Should be able to support the following commands:
// - tree
// - cwp
// - mv
type Tree struct {
	client       *spotify.Client         // Reference to Spotify client for mounting purposes
	cwp          *spotify.SimplePlaylist // Current working playlist
	children     []Node                  // Subdirectories
	num_children int
}

// Fetches playlists within entire user library
// Public scoping because it'd be pretty useful
func FetchPlaylists() []spotify.SimplePlaylist {
	var ret []spotify.SimplePlaylist
	offset := 0
	limit := 20

	for {
		playlists, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{Offset: &offset, Limit: &limit})
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

// Constructs a folder by doing a depth-first search on a list of playlists
// Because the API returns folder content in linear order, we can treat this as DFS
func constructFolder(playlists []spotify.SimplePlaylist, dirname string, index int, folders map[string]string) (Node, int) {
	var children []Node
	iter := 0

	for i := index; i < len(playlists) && folders[string(playlists[i].URI)] == dirname; i++ {
		playlist := playlists[i]
		children = append(children, constructPlaylist(playlist.Name, playlist.ID))
		iter++
	}
	node := Node{
		Name:         dirname,
		Format:       "folder",
		Children:     children,
		Id:           "",
		Num_children: len(children),
	}
	return node, iter
}

// Constructs a Node of format "Playlist"
// I don't think we need to store the tracks as children;
// instead, we should save the ID so that we can retrieve
// the tracks on demand
func constructPlaylist(name string, playlist_id spotify.ID) Node {
	return Node{
		Name:         name,
		Format:       "playlist",
		Children:     nil,
		Id:           playlist_id,
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
		if node.Format == "folder" {
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
		client:   client,
		cwp:      nil, // figure out what to do for root
		children: nodes,
	}
}
