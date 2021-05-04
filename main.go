package main

// Imported packages
import (
	"fmt"
	"log"

	"os"
	"bufio"
	"strings"

	"errors"

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


func execInput(input string) error {
	// Tokenize input by removing newline 
	input = strings.TrimSuffix(input, "\n")

	// Handle execution of input
	// Switch cases for subcommands
	// Need to figure out how to make HTTP requests from commands module, or do it from here via commands
	// Somehow need to pass user globally; gotta figure that out
	switch input {
		case "hello":
			cmd.HelloInit(os.Args[1:])
		case "exit":
			os.Exit(1)
		default:
			return errors.New("Oopsies")
	}
	return nil
}

// Executes the interactive shell
func startShell() {
	// Read in standard inputs
	reader := bufio.NewReader(os.Stdin)
	// Shell should loop infinitely unless sent SIGINT is raised
	for {
		fmt.Print("> ")
		// Read keyboard inputs
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

func main() {
	// Set authentication details
	client_id, client_secret := fetchCredentials(".credentials.json")
	//fmt.Printf("client_id=%s\nclient_secret=%s", client_id, client_secret)
	auth.SetAuthInfo(client_id, client_secret)

	// Start HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for: ", r.URL.String())
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
	fmt.Printf("You are logged in as: %s\n", user.ID) // Need to figure out how to cache this
	
	// At this point, create directory structure

	// Run shell
	startShell()
}