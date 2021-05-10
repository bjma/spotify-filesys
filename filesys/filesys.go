package filesys

import (
    "bytes"
    "fmt"

    "github.com/bjma/spotify-filesys/api"
)

var (
    dir_tree  *Tree = nil
    tree_flag       = false
)

// When things get bigger, we should actually build a clean filesystem
// This system should act as the layout for our commands
func ParseDir(T *Tree) []*Node {
    var ret []*Node
    return append(ret, &Node{
        Name:     ".",
        Format:   "root",
        Children: T.Children,
        Id:       "",

        Num_children: len(T.Children),
    })
}

// Reads array of Nodes and returns its children and max_depth possible
func ScanDir(dir []*Node, dname string) ([]*Node, int) {
    // If root, return children of root
    // If not root, return child
    // Essentially we just want to return subarrays recursively
    return nil, -1
}

func PrintTree(T *Tree, depth int) {
    dir_tree = T

    var w_buf bytes.Buffer
    treeRecurse(T.Children, ".", depth, &w_buf, 0, false)
    fmt.Println(w_buf.String())
}

// Debuggin Purposes
func treeIter(T *Tree, buf *bytes.Buffer, depth int) {
    tree := T.Children

    buf.WriteString(".\n")
    for i, node := range tree {
        if i == T.Num_children-1 {
            //fmt.Printf("%s%s\n", "└──", node.Name)
            buf.WriteString("└──")
        } else {
            //fmt.Printf("%s%s\n", "├──", node.Name)
            buf.WriteString("├──")
        }
        buf.WriteString(node.Name + "\n")

        if node.Format == "folder" || node.Format == "artist" {
            for j, child := range node.Children {
                if i != T.Num_children {
                    buf.WriteString("|  ")
                } else {
                    buf.WriteString("   ")
                }
                if j == node.Num_children-1 {
                    buf.WriteString("└──")
                } else {
                    buf.WriteString("├──")
                }
                buf.WriteString(child.Name + "\n")

                // Print up to depth tracks
                tracks := api.FetchTracks(child.Id, depth, child.Format)
                for k, track := range tracks {
                    if j != node.Num_children-1 {
                        buf.WriteString("   |  ")
                    } else {
                        buf.WriteString("      ")
                    }
                    if k == len(tracks)-1 {
                        buf.WriteString("└──")
                    } else {
                        buf.WriteString("├──")
                    }
                    buf.WriteString(track.Name + "\n")
                }
            }
        }
    }
}

func indent(level int, offset int) string {
    var ret string

    num_spaces := (level + offset) * 3
    for i := 0; i < num_spaces; i++ {
        ret += " "
    }
    return ret
}

// Prints tree in preorder traversal by writing each level into 
// a write buffer. The flag is used to signal whether or not the
// current node is a leaf or not.
func treeRecurse(T []*Node, dirname string, depth int, w_buf *bytes.Buffer, level int, is_leaf bool) {
    if level > 0 {
        w_buf.WriteString(indent(level, 0))
        //if is_leaf {
        //    w_buf.WriteString("└──")
        //} else {
        //    w_buf.WriteString("├──")
        //}
    }
    w_buf.WriteString(dirname)
    w_buf.WriteString("\n")

    // For each current level, recurse on their children.
    // If current child is a playlist or album, print up
    // to `depth` tracks.
    for _, child := range T {
        if child != nil {
            treeRecurse(child.Children, child.Name, depth, w_buf, level+1, child.Is_leaf)

            if child.Children == nil {
                tracks := api.FetchTracks(child.Id, depth, child.Format)
                for _, track := range tracks {
                    //w_buf.WriteString(indent(level, 1))
                    //if i == len(tracks)-1 {
                    //    w_buf.WriteString("└──")
                    //} else {
                    //    w_buf.WriteString("├──")
                    //}
                    w_buf.WriteString(indent(level, 2))
                    w_buf.WriteString(track.Name)
                    w_buf.WriteString("\n")
                }
            }
        }
    }
}
