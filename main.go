package main

import (
    "fmt"
    "os"
    "errors"
    "bufio"
    "strings"
    "encoding/json"
    "io/ioutil"

    "github.com/bjma/spotify-filesys/api"
    "github.com/bjma/spotify-filesys/filesys" // Filesystem
    "github.com/bjma/spotify-filesys/cmd"     // Subcommands
)
// Global reference to directory tree
var tree *filesys.Tree

// Parses config file and returns authentication details 
// and a map representation of user's folder hierarchy
// Source: https://golangr.com/read-json-file/
func readConfig(filename string) (string, string, map[string]string) {
    // Playlists in config.folders.children field
    type Playlist struct {
        Type string `json: "type"`
        Uri  string `json: "uri"`
    }
    // Folder objects in config.folders field
    type Folder struct {
        Uri      string     `json: "uri"`
        Type     string     `json: "type"`
        Children []Playlist `json: "children"`
        Name     string     `json: "name"`
    }
    // Config file struct
    type Config struct {
        Client_id     string   `json: "client_id"`
        Client_secret string   `json: "client_secret"`
        Folders       []Folder `json: "folders"`
    }

    file, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }

    buffer := Config{}
    if err = json.Unmarshal(file, &buffer); err != nil {
        panic(err)
    }

    // Maps playlist URIs to owner name for parsing folders
    folders := make(map[string]string)
    for _, folder := range buffer.Folders {
        for _, playlist := range folder.Children {
            folders[playlist.Uri] = folder.Name
        }
    }

    return buffer.Client_id, buffer.Client_secret, folders
}

// Execute input commands within shell
func execInput(input string) error {
    // Tokenize input against whitespace and remove newline
    input = strings.TrimSuffix(input, "\n")
    args := strings.Split(input, " ")
    input = args[0]

    switch input {
    case "hello":
        cmd.HelloInit(args[1:])
    case "whoami":
        fmt.Printf("%s (%s)\n", api.GetUserName(), api.GetUserID())
    case "tree":
        cmd.TreeInit(tree, args[1:])
    case "exit":
        os.Exit(1)
    default:
        return errors.New("spfs: command not found: " + input)
    }
    return nil
}

func main() {
    // Set authentication details
    client_id, client_secret, folders := readConfig(".config.json")

    // Handle authentication
    err := api.HandleAuth(client_id, client_secret)
    if err != nil {
        panic(err)
    }
    // Construct filesystem and begin interactive shell
    tree = filesys.BuildTree(folders, "user")

    // Shell should loop infinitely unless sent SIGINT is raised or `exit` is executed
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("> ")
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }

        // Execute commands
        if err = execInput(input); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}
