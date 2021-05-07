package main

import (
	"fmt"
	"log"

	"bufio"
	"os"
	"strings"

	"errors"

	"net/http"

	"encoding/json"
	"io/ioutil"

	// Libraries
	"github.com/zmb3/spotify"

	// Modules
	"github.com/bjma/spotify-filesys/filesys" // Filesystem
	"github.com/bjma/spotify-filesys/cmd"     // Subcommands
)

const redirectURI = "http://localhost:8080/callback"
// Global reference to directory tree
var tree *filesys.Tree

// Authentication details
var (
	auth = spotify.NewAuthenticator(redirectURI,
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserTopRead,
		spotify.ScopePlaylistReadPrivate,
	)
	ch                     = make(chan *spotify.Client)
	state                  = "abc123"
	client *spotify.Client = nil
)

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

// Completes authentication by verifying credentials
// Source: https://github.com/zmb3/spotify/blob/master/examples/authenticate/authcode/authenticate.go
func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s", st, state)
	}
	// Use token to get authenticated client
	client := auth.NewClient(token)
	fmt.Fprintf(w, "Login completed!")
	ch <- &client

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You are logged in as: %s\n\n", user.ID)
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
		user, _ := client.CurrentUser()
		fmt.Println(user.ID)
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
	auth.SetAuthInfo(client_id, client_secret)

	// Start HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for: ", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Printf("Please log in to Spotify by visiting the following page in your browser:%s\n\n", url)

	// Wait for auth to complete
	client = <- ch

	// Construct filesystem and begin interactive shell
	tree = filesys.BuildTree(client, folders, "user")

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
