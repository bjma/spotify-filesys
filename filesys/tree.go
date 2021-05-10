package filesys

import (
    "fmt"
    "github.com/zmb3/spotify"

    "github.com/bjma/spotify-filesys/api"
)

// TODO:
// Maybe clean up the actual design of the tree
// i.e., design things with more intention
// also, maybe use linked lists if possible, over arrays?
// Use thsi as example:
// https://github.com/rfbergeron/bompiler/blob/master/astree.c

// Abstraction for tracks, artists, albums, playlists, and folders
type Node struct {
    Name         string
    Format       string     // Tracks, Artists, Playlists, Folders
    Children     []*Node    // If artist or folders, Children != nil to store playlists and folders
    Id           spotify.ID // If album or playlist, Id != nil to retrieve tracks by ID
    Num_children int
    Is_leaf      bool // Makes it easier for us to print trees
}

// Essentially the data structure that represents the filesystem;
// Navigation and all other stuff will be done on a Tree (see `fs.go`)
// NOTE: Filesystem should only be done with user library (can maybe try to add saved albums)
type Tree struct {
    Client       *spotify.Client // Reference to Spotify client for mounting purposes
    Cwp          *Node           // Current working playlist (ls will display children of Cwp)
    Children     []*Node         // Subdirectories
    Num_children int
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
    // Really jank way to guarantee that last directory entry is a leaf
    children[len(children)-1].Is_leaf = true

    node := &Node{
        Name:         dirname,
        Format:       "folder",
        Children:     children,
        Id:           "",
        Num_children: len(children),
        Is_leaf:      false,
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
        Is_leaf:      is_leaf,
    }
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
            Is_leaf:      i == n-1,
        })
        i++
    }

    return nodes
}

func constructUserTree(folders map[string]string) []*Node {
    var nodes []*Node

    lib := api.FetchPlaylists()

    i := 0
    n := len(lib)
    // Iterate through user's entire playlist library and initialize data needed to generate directory tree
    for i < n {
        uri := string(lib[i].URI)

        // If current playlist is in a folder, parse it as a folder, then append to playlist
        // Else, just append to nodes as playlist
        if folders[uri] != "" {
            node, iter := newFolder(lib, folders[uri], i, folders)
            nodes = append(nodes, node)
            i += iter
        } else {
            nodes = append(nodes, newPlaylist(lib[i].Name, lib[i].ID, i == n-1))
            i++
        }
    }
    // This is so stupid lol
    nodes[len(nodes)-1].Is_leaf = true
    return nodes
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

// Runs DFS on given directory and returns the destination
// This is kind of jank when it goes down to further levels;
// just returns an empty node (fix later)
func GetNodeByName(dir []*Node, dirname string) *Node {
    var ret Node
    
    if dir == nil {
        return nil
    }

    for _, subdir := range dir {
        if subdir.Name == dirname {
            ret = *subdir
            fmt.Printf("is %s a leaf? %t\n", ret.Name, ret.Is_leaf)
            break
        }
        GetNodeByName(subdir.Children, dirname)
    }
    return &ret
}

// Builds a directory tree
func BuildTree(f map[string]string, opt string) *Tree {
    nodes := parseLibrary(opt, f)
    root := Node{
        Name:         ".",
        Format:       "root",
        Children:     nodes,
        Id:           "",
        Num_children: len(nodes),
    }

    return &Tree{
        Cwp:          &root, // figure out what to do for root
        Children:     nodes,
        Num_children: len(nodes),
    }
}
