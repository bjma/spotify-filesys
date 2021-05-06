package cmd

import (
	"flag"
	_ "flag"
	_ "fmt"

	"github.com/bjma/spotify-filesys/filesys"
)

var treeCommand *flag.FlagSet
var T *filesys.Tree

func TreeInit(tree *filesys.Tree, arr []string) {
	treeCommand = flag.NewFlagSet("tree", flag.ExitOnError)
	T = tree

	// Flags here
	artistTreePtr := treeCommand.Bool("artists", false, "Display followed artists and their albums/tracks")
	albumTreePtr := treeCommand.Bool("albums", false, "Display saved albums")

	treeCommand.Parse(arr)

	if treeCommand.Parsed() {
		// 0 - user lib; 1 - artist; 2 - saved albums
		if *artistTreePtr {
			TreeExecute(1)
		} else if *albumTreePtr {
			TreeExecute(2)
		} else {
			TreeExecute(0)
		}
	}
}

func TreeExecute(opt int) {
	switch opt {
	case 0:
		filesys.PrintTree(T, false)
	case 1:
		// Generate artist tree
		artists := filesys.BuildTree(T.Client, nil, 1)
		filesys.PrintTree(artists, false)
		break
	case 2:
		// Generate album tree
		break
	}
}