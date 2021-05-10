package cmd

import (
    "flag"
    "fmt"

    "github.com/bjma/spotify-filesys/filesys"
)

var treeCommand *flag.FlagSet
var T *filesys.Tree

func TreeInit(tree *filesys.Tree, arr []string) {
    treeCommand = flag.NewFlagSet("tree", flag.ExitOnError)
    T = tree

    // Flags here
    artistTree := treeCommand.Bool("artists", false, "Display followed artists and their albums/tracks")
    albumTree := treeCommand.Bool("albums", false, "Display saved albums")
    trackLimit := treeCommand.Int("depth", 0, "Show up to DEPTH tracks")

    treeCommand.Parse(arr)

    if treeCommand.Parsed() {
        if *artistTree {
            TreeExecute("artists", *trackLimit)
        } else if *albumTree {
            TreeExecute("albums", *trackLimit)
        } else {
            TreeExecute("user", *trackLimit)
        }
    }
}

func TreeExecute(opt string, depth int) (err error) {
    var t *filesys.Tree

    // Some rule enforcement
    if depth > 5 || depth < 0 {
        err = fmt.Errorf("spfs: error: depth limit exceeded")
        fmt.Println(err)
        return err
    }

    switch opt {
    case "user":
        t = T
    case "artists":
        // Generate artist tree
        t = filesys.BuildTree(nil, "artists")
    case "albums":
        // Generate album tree
        break
    }
    filesys.PrintTree(t, depth)
    return
}