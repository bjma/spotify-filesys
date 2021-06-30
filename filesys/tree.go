package filesys

import (
	"bytes"
	"fmt"
	"github.com/zmb3/spotify"

	"github.com/bjma/spotify-filesys/api"
)

var (
	flag = false
)

type Node struct {
	Name         string
	Format       string     // Tracks, Artists, Playlists, Folders
	Children     []*Node    // If artist or folders, Children != nil to store playlists and folders
	Id           spotify.ID // If album or playlist, Id != nil to retrieve tracks by ID
	Num_children int
}

type Tree struct {
	Client       *spotify.Client // Reference to Spotify client for mounting purposes
	Cwp          *Node           // Current working playlist (ls will display children of Cwp)
	Dir         []*Node          // Library tree
	Num_children int
}

// Builds a directory tree
func BuildTree(f map[string]string, opt string) *Tree {
    var libTree []*Node

    lib :=  parseLibrary(opt, f)
	root := Node{
		Name:         ".",
		Format:       "root",
		Children:     lib,
		Id:           "",
		Num_children: len(lib),
    }
    libTree = append(libTree, &root)

	return &Tree{
		Cwp:          &root, // figure out what to do for root
		Dir:     libTree,
		Num_children: len(libTree),
	}
}

// `tree` command
func PrintTree(T *Tree, depth int) {
	dir_tree = T

	var w_buf bytes.Buffer
    //printTreeRecurse(T.Children, ".", depth, &w_buf, 0, false)
    printTreeTest(T.Dir, &w_buf)
	fmt.Println(w_buf.String())
}

// Parses entire user library for building filesystem
func parseLibrary(opt string, folders map[string]string) []*Node {
	var nodes []*Node

	switch opt {
	case "user":
		nodes = constructUserTree(folders)
	case "artists":
		nodes = constructArtistTree()
	case "albums":
		break
	}
	return nodes
}

// NOTE: Might just have a single function to construct trees
// so that we don't clutter code
func constructUserTree(folders map[string]string) []*Node {
	var nodes []*Node

	lib := api.FetchPlaylists()

	i := 0
	n := len(lib)

	// Iterate through user's entire playlist library and initialize data needed to
	// generate directory tree
	for i < n {
		uri := string(lib[i].URI)

		if folders[uri] != "" {
			node, iter := newFolder(lib, folders[uri], i, folders)
			nodes = append(nodes, node)
			i += iter
		} else {
			nodes = append(nodes, newPlaylist(lib[i].Name, lib[i].ID, i == n-1))
			i++
		}
	}

	return nodes
}

func constructArtistTree() []*Node {
	var nodes []*Node

	lib := api.FetchArtists()

	i := 0
	n := len(lib)
	for i < n {
		// Get list of artist albums
		albums := api.FetchAlbums(lib[i].ID, "artist")
		var children []*Node

		for _, album := range albums {
			children = append(children, &Node{
				Name:         album.Name,
				Format:       "album",
				Children:     nil,
				Id:           album.ID,
				Num_children: 0,
			})
		}

		nodes = append(nodes, &Node{
			Name:         lib[i].Name,
			Format:       "artist",
			Children:     children,
			Id:           "",
			Num_children: len(children),
		})
		i++
	}

	return nodes
}

// Initializes a folder Node by doing a depth-first search on a list of playlists
// Returns offset for index to avoid unnecessary searches/parsing
func newFolder(playlists []spotify.SimplePlaylist, dirname string, index int, folders map[string]string) (*Node, int) {
	var children []*Node
	iter := 0

	n := len(playlists)
	for i := index; i < n && folders[string(playlists[i].URI)] == dirname; i++ {
		playlist := playlists[i]
		children = append(children, newPlaylist(playlist.Name, playlist.ID, false))
		iter++
	}

	node := &Node{
		Name:         dirname,
		Format:       "folder",
		Children:     children,
		Id:           "",
		Num_children: len(children),
	}
	return node, iter
}

// Initializes a Node of format "Playlist"
func newPlaylist(name string, playlist_id spotify.ID, is_leaf bool) *Node {
	return &Node{
		Name:         name,
		Format:       "playlist",
		Children:     nil,
		Id:           playlist_id,
		Num_children: 0,
	}
}

// Helper function for treeRecurse that prints up to `depth` tracks
// Also I want to get a better naming for this function
func printTracks(id spotify.ID, format string, w_buf *bytes.Buffer, depth int, level int) {
	var indent string
	num_spaces := level * 3
	for i := 0; i < num_spaces; i++ {
		indent += " "
	}
	tracks := api.FetchTracks(id, depth, format)
	for i, track := range tracks {
		if flag {
			w_buf.WriteString("|")
		}
		w_buf.WriteString(indent)
		if i == len(tracks)-1 {
			w_buf.WriteString("└──")
		} else {
			w_buf.WriteString("├──")
		}

		w_buf.WriteString(track.Name)
		w_buf.WriteString("\n")
	}
}

func printTreeTest(tree []*Node, w_buf *bytes.Buffer) {
    if tree == nil {
        return
    }

    for _, node := range tree {
        fmt.Printf("%s\n", node.Name)
        printTreeTest(node.Children, w_buf)
    }
}