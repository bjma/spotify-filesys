package api

import (
    "github.com/zmb3/spotify"
    "fmt"
    "log"
    "net/http"
)

const redirectURI = "http://localhost:8080/callback"

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

// Completes authentication by verifying credentials
// Source: https://github.com/zmb3/spotify/blob/master/examples/authenticate/authcode/authenticate.go
func CompleteAuth(w http.ResponseWriter, r *http.Request) {
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
    fmt.Printf("You are logged in as: %s\n\n", user.DisplayName)
}

func HandleAuth(client_id string, client_secret string) error {
    auth.SetAuthInfo(client_id, client_secret)

    // Start HTTP server
    http.HandleFunc("/callback", CompleteAuth)
    go http.ListenAndServe(":8080", nil)

    url := auth.AuthURL(state)
    fmt.Printf("Please log in to Spotify by visiting the following page in your browser:%s\n\n", url)

    // Wait for auth to complete
    client = <- ch
    return nil // No reason for raising errors yet
}