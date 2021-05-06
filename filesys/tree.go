package filesys

import (
	"fmt"

	// Libraries
	"github.com/zmb3/spotify"
)

var (
	client *spotify.Client = nil
	tree_level = 0 // For tree recursion
)

/*
 * Abstraction for tracks, playlists, and folders
 * If Format is type:
 * - Folder/Artist: childrens list contains Nodes of Format `playlist`; no ID is assigned
 * - Playlist/Album: childrens list is empty; instead, ID is assigned to retrieve tracks on request
 */
type Node struct {
	Name         string     // Playlist or folder name
	Format       string     // {"folder", "playlist", "album", "artist", "track"}
	Children     []Node     // Children of node
	Id           spotify.ID // If playlist, we save its ID for access to tracks
	Num_children int
}

type Tree struct {
	Client       *spotify.Client         // Reference to Spotify client for mounting purposes
	Cwp          *spotify.SimplePlaylist // Current working playlist
	Children     []Node                  // Subdirectories
	Num_children int
}

// Fetches all playlists in user library
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

// Fetches user's Top 10 Artists (all followed artists too slow)
func FetchArtists() []spotify.SimpleArtist {
	var ret []spotify.SimpleArtist
	limit := 10

	artists, err := client.CurrentUsersTopArtistsOpt(&spotify.Options{Limit: &limit})
	if err != nil || artists == nil {
		panic(err)
	}

	for _, artist := range artists.Artists {
		ret = append(ret, artist.SimpleArtist)
	}
	return ret
}

// Fetches entire list of albums; 
// If `opt` = 1, return artist discography
// if `opt` = 2, return saved albums
func FetchAlbums(id spotify.ID, opt int) []spotify.SimpleAlbum {
	var ret []spotify.SimpleAlbum
	offset := 0
	var limit int

	switch opt {
	case 1:
		limit = 5

		albums, err := client.GetArtistAlbumsOpt(id, &spotify.Options{Offset: &offset, Limit: &limit})
		if err != nil || albums == nil || len(albums.Albums) < 1 {
			panic(err)
		}

		for _, album := range albums.Albums {
			ret = append(ret, album)
		}
	case 2:
		limit = 20

		for {
			albums, err := client.CurrentUsersAlbumsOpt(&spotify.Options{Offset: &offset, Limit: &limit})
			if err != nil || albums == nil || len(albums.Albums) < 1 {
				break
			}
	
			for _, album := range albums.Albums {
				ret = append(ret, album.SimpleAlbum)
			}
			offset += limit
		}
	}
	return ret
}

// Initializes a folder Node by doing a depth-first search on a list of playlists
// Returns offset for index to avoid unnecessary searches/parsing
func newFolder(playlists []spotify.SimplePlaylist, dirname string, index int, folders map[string]string) (Node, int) {
	var children []Node
	iter := 0

	for i := index; i < len(playlists) && folders[string(playlists[i].URI)] == dirname; i++ {
		playlist := playlists[i]
		children = append(children, newPlaylist(playlist.Name, playlist.ID))
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

// Initializes a Node of format "Playlist"
func newPlaylist(name string, playlist_id spotify.ID) Node {
	return Node{
		Name:         name,
		Format:       "playlist",
		Children:     nil,
		Id:           playlist_id,
		Num_children: 0,
	}
}

func constructArtistTree() []Node {
	var nodes []Node 

	lib := FetchArtists()

	i := 0
	for i < len(lib) {
		// Get list of artist albums
		albums := FetchAlbums(lib[i].ID, 1)
		var children []Node

		for _, album := range albums {
			children = append(children, Node{
				Name: album.Name,
				Format: "album",
				Children: nil,
				Id: album.ID,
				Num_children: 0,
			})
		}

		nodes = append(nodes, Node{
			Name: lib[i].Name,
			Format: "artist",
			Children: children,
			Id: "",
			Num_children: len(children),
		})
		i++
	}

	return nodes
}

func constructUserTree(folders map[string]string) []Node {
	var nodes []Node
	
	lib := FetchPlaylists()

	i := 0
	// Iterate through user's entire playlist library and initialize data needed to generate directory tree
	for i < len(lib) {
		uri := string(lib[i].URI)

		// If current playlist is in a folder, parse it as a folder, then append to playlist
		// Else, just append to nodes as playlist
		if folders[uri] != "" {
			node, iter := newFolder(lib, folders[uri], i, folders)
			nodes = append(nodes, node)
			i += iter
		} else {
			nodes = append(nodes, newPlaylist(lib[i].Name, lib[i].ID))
			i++
		}
	}
	return nodes
}

// Parses entire user library for building filesystem
// Need to also figure out a way to do this for Saved Albums and user's Top Artists
func parseLibrary(opt int, folders map[string]string) []Node {
	var nodes []Node
	
	switch opt {
	case 0:
		nodes = constructUserTree(folders)
	case 1:
		nodes = constructArtistTree()
	case 2:
		break
	}
	return nodes
}

// Debugging purposes
// Should actually try to do this in a DFS like way
func PrintTree(t *Tree, flag bool) {
	tree := t.Children

	fmt.Println(".")
	for i, node := range tree {
		if i == t.Num_children - 1 {
			fmt.Printf("%s%s\n", "└──", node.Name)
		} else {
			fmt.Printf("%s%s\n", "├──", node.Name)
		}
		if node.Format == "folder" || node.Format == "artist" { 
			for j, child := range node.Children {
				if i != t.Num_children - 1 {
					fmt.Printf("|  ")
				} else {
					fmt.Printf("   ")
				}
				if j == node.Num_children - 1 {
					fmt.Printf("%s%s\n", "└──", child.Name)
				} else {
					fmt.Printf("%s%s\n", "├──", child.Name)
				}
			}
		}
	}
}

// Builds a directory tree
// Options:
// - 0 (user library)
// - 1 (followed artists)
// - 2 (saved albums)
func BuildTree(c *spotify.Client, f map[string]string, opt int) *Tree {
	client = c

	nodes := parseLibrary(opt, f)

	return &Tree{
		Client:   client,
		Cwp:      nil, // figure out what to do for root
		Children: nodes,
		Num_children: len(nodes),
	}
}
