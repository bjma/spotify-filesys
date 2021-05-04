package main

// Imported packages
import (
	"fmt"
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
	// Libraries
	"github.com/zmb3/spotify"

	// Modules
	"github.com/bjma/spotify-filesys/cmd" // Subcommands
	_ "github.com/bjma/spotify-filesys/tests" // Unit testing
)

// Authentication details
const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	ch = make(chan *spotify.Client)
	state = "abc123"
)

// Parses JSON file for credentials and returns client ID and secret
// Source: https://golangr.com/read-json-file/
func fetchCredentials(filename string) (id, secret string) {
	// I guess structures make parsing JSONs easier (but leaves lots of trash)
	type Credentials struct {
		Client_id string `json: "client_id"`
		Client_secret string `json: "client_secret"`
	}
	
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	
	buffer := Credentials{}
	// Unmarshal JSON
	if err = json.Unmarshal(file, &buffer); err != nil {
		panic(err)
	}

	return buffer.Client_id, buffer.Client_secret
}

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
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: Not enough arguments.")
		// Print usage errors; prob can use a package for that
		os.Exit(1)
	}

	// Set authentication details
	client_id, client_secret := fetchCredentials(".credentials.json")
	//fmt.Printf("client_id=%s\nclient_secret=%s", client_id, client_secret)
	auth.SetAuthInfo(client_id, client_secret)

	// Start HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// Wait for auth to complete (def go back and revisit these)
	client := <-ch

	// Use client to make API calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	
	// At this point, create directory structure

	// Switch cases for subcommands
	// Need to figure out how to make HTTP requests from commands module, or do it from here via commands
	// Somehow need to pass user globally; gotta figure that out
	switch os.Args[1] {
		// Testing
		case "hello":
			cmd.HelloInit(os.Args[2:])
		default:
			os.Exit(1)
	}
}