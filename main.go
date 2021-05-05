package main

// Imported packages
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
	"github.com/bjma/spotify-filesys/cmd"     // Subcommands
	"github.com/bjma/spotify-filesys/filesys" // Filesystem
)

// Authentication details
const redirectURI = "http://localhost:8080/callback"

var (
	// Idk if i like looking at this
	auth = spotify.NewAuthenticator(redirectURI,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistReadPrivate,
	)
	ch                     = make(chan *spotify.Client)
	state                  = "abc123"
	client *spotify.Client = nil
)

// Global reference to directory tree
var tree *filesys.Tree

// Parses config file and returns a map representation of JSON content
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

	// Write JSON into buffer
	buffer := Config{}
	if err = json.Unmarshal(file, &buffer); err != nil {
		panic(err)
	}

	// Additional stuff for folder parsing
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
	fmt.Fprintf(w, "Login complete!")
	ch <- &client // read up on the arrows

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("You are logged in as: %s\n\n", user.ID) // Need to figure out how to cache this
}

// Execute input commands within shell
func execInput(input string) error {
	// Tokenize input against whitespace and remove newline
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")
	input = args[0]

	// Handle execution of input
	// Switch cases for subcommands
	switch input {
	case "hello": // Testing purposes
		cmd.HelloInit(args[1:])
	case "whoami":
		user, _ := client.CurrentUser()
		fmt.Println(user.ID)
	case "tree":
		filesys.PrintTree(tree)
	case "exit":
		os.Exit(1)
	default:
		return errors.New("computer said to tell u that ur fucking stupid")
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
	tree = filesys.BuildTree(client, folders)

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
