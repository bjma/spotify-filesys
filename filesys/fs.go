package filesys

import (
	"bytes"
	"fmt"
)

var (
	tree_flag = false
)

// When things get bigger, we should actually build a clean filesystem
// This system should act as the layout for our commands

// Runs DFS on given directory and returns the destination
func DFS(dir []*Node, dirname string) *Node {
	if dir == nil {
		return nil
	}

	for _, subdir := range dir {
		if subdir.Name == dirname {
			return subdir
		}
		DFS(subdir.Children, dirname)
	}
	return nil
}

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
				tracks := FetchTracks(child.Id, depth, child.Format)
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

// TODO:
// * Make recursive with indentation buffer so we don't need to spaghetti code this shit
// * Write tree into buffer so it doesn't print line-by-line
// * Add track.Artists too
func 
treeRecurse(T []*Node, depth int, w_buf *bytes.Buffer, indent_buf *bytes.Buffer, level int) string {
	// Algorithm
	// 1. ScanDir -> if max_depth < 0, then return; else, recurse on children
	// 2. For each child in current level -> 
	// 3. Write indent buffer and current child name into buffer
	// 4. If depth flag is raised (i.e. depth > 0), then write up to depth tracks with added indentation buffer into write buffer
	//    -- NOTE: Should only do this if current node is a playlist/album
	// 5. Recurse

	// Note that we should ALSO keep track of the format of the current Node
	// i.e. if Node is "Album" or "Playlist, " print up to DEPTH tracks
	// If Node is "Artist" or "Folder", just print children
	// Indentation buffer should be able to aid with formatting easily
	if T == nil {
		return ""
	}

	if level > 1 {
		indent_buf.WriteString("   ")
	}

	orig_buf := indent_buf
	for i := 0; i < len(T); i++ {
		//fmt.Printf("%s\n", T[i].Name)
		buf := indent_buf 
		buf.WriteString("└──")
		buf.WriteString(T[i].Name)
		buf.WriteString("\n")

		w_buf.WriteString(buf.String())

		treeRecurse(T[i].Children, depth, w_buf, indent_buf, level)
		indent_buf = orig_buf
	}
	return ""
}

func PrintTree(T *Tree, depth int) {
	var w_buf bytes.Buffer
	//var indent_buf bytes.Buffer
	treeIter(T, &w_buf, depth)
	//dir := ParseDir(T)
	//treeRecurse(dir, depth, &w_buf, &indent_buf, 1)
	fmt.Println(w_buf.String())
}
