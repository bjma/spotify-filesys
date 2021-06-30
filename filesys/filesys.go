package filesys

import (
	_ "fmt"
)

var (
	dir_tree *Tree = nil
)

// Reads directory `dirname` and creates
// an array of the total number of elements
// 
//func ScanDir(dirname string, dir []*Node) {}
