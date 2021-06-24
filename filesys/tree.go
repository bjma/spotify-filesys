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

// `tree` command
func PrintTree(T *Tree, depth int) {
    dir_tree = T

    var w_buf bytes.Buffer
    printTreeRecurse(T.Children, ".", depth, &w_buf, 0, false)
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

// Prints tree in preorder traversal by writing each level into
// a write buffer. The flag is used to signal whether or not the
// current node is a leaf or not.
//
// NOTE: Stop working on treeRecurse for now. Instead, begin implementing ScanDir
// so that we can accurately determine our position relative to current path level.
// Also I kind of want a better naming for this function.
func printTreeRecurse(T []*Node, dirname string, depth int, w_buf *bytes.Buffer, level int, is_leaf bool) {
    if level > 1 {
        // If still printing children, then print VER
        // Else, print extra space
        if flag {
            w_buf.WriteString("|")
        } else {
            w_buf.WriteString(" ")
        }

        num_spaces := level * 3 - 3
        for i := 0; i < num_spaces-1; i++ {
            w_buf.WriteString(" ")
        }
    }
    if level > 0 {
        // is_leaf is a temporary fix. We should implement ScanDir to
        // tell us whether or not the current item is the last entry
        // in its parent. If so, every time when we return to the 
        // next recursive level, we should also be able to tell
        // if the current level is the last in its subdirectory.
        // Or rather, is_leaf should only be used to determine 
        // between ELB and TEE; raising/putting down the flag
        // should then be decided by some other metric (like ScanDir)
        if is_leaf {
            w_buf.WriteString("└──")
        } else {
            // If we return from the last subdirectory back to another subdirectory, put flag down
            if flag {
                flag = false
            }
            w_buf.WriteString("├──")
        }
    }
    w_buf.WriteString(dirname)
    w_buf.WriteString("\n")

    // For each current level, recurse on their children.
    // If current child is a playlist or album, print up
    // to `depth` tracks.
    for _, child := range T {
        if child != nil {
            printTreeRecurse(child.Children, child.Name, depth, w_buf, level+1, child.Is_leaf)

            if child.Children == nil {
                printTracks(child.Id, child.Format, w_buf, depth, level+1)
            }
        }
    }
}
