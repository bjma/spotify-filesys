package filesys

import (
    _"fmt"
)

var (
    dir_tree  *Tree = nil
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

// Finds the position of dirname relative to its directory level
// If it's the last element, return -1
// Else, return the max depth remaining relative to current level
// Requirements: BFS
func ScanDir(dir []*Node, dname string) ([]*Node, int) {
    // We'd want to do this in a sort of BFS like fashion;
    // Rather, find the level that `dname` belongs to,
    // get the position of `dname` relative to that level (i.e. 5th child),
    // and just return the remaining children left in the level
    return nil, -1
}

// Runs DFS on given directory and returns the destination
// NOTE: There's a problem where the node pointer returned 
// is corrupted; that is, within the function scope of DFS,
// the data exists, but once returned the object is empty.
// This is most likely due to how we're dealing with pointers
// (i.e. lack of understanding of Golang), but also might
// be how we're recursing. 
// Look into the recursive calls via *unit testing*,
// maybe we need to approach this from a different angle.
func DFS(dir []*Node, dirname string) *Node {
    var ret *Node
    
    if dir == nil {
        return nil
    }

    for _, subdir := range dir {
        if subdir.Name == dirname {
            ret = subdir
            break
        }
        DFS(subdir.Children, dirname)
    }
    return ret // is returned nil most of the time
}