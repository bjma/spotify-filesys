package api

// Package to handle all requests to the Spotify API
// Deals with handling authentication on the client side,
// as well as serving fetch requests from each program.

import (
    "github.com/zmb3/spotify"
    _ "fmt"
)


// Fetches user's Top 10 Artists (all followed artists too slow)
func FetchArtists() []spotify.SimpleArtist {
    var ret []spotify.SimpleArtist
    limit := 10

    artists, err := client.CurrentUsersTopArtistsOpt(&spotify.Options{Limit: &limit})
    if err != nil || artists == nil {
        panic(err)
    }

    for _, artist := range artists.Artists {
        ret = append(ret, artist.SimpleArtist)
    }
    return ret
}

// Fetches entire list of albums;
// If `opt` = 1, return artist discography
// if `opt` = 2, return saved albums
func FetchAlbums(id spotify.ID, opt string) []spotify.SimpleAlbum {
    var ret []spotify.SimpleAlbum
    offset := 0
    var limit int

    switch opt {
    case "artist":
        limit = 5

        albums, err := client.GetArtistAlbumsOpt(id, &spotify.Options{Offset: &offset, Limit: &limit})
        if err != nil || albums == nil || len(albums.Albums) < 1 {
            panic(err)
        }

        for _, album := range albums.Albums {
            ret = append(ret, album)
        }
    case "user":
        limit = 20

        for {
            albums, err := client.CurrentUsersAlbumsOpt(&spotify.Options{Offset: &offset, Limit: &limit})
            if err != nil || albums == nil || len(albums.Albums) < 1 {
                break
            }

            for _, album := range albums.Albums {
                ret = append(ret, album.SimpleAlbum)
            }
            offset += limit
        }
    }
    return ret
}

// Fetches all playlists in user library
func FetchPlaylists() []spotify.SimplePlaylist {
    var ret []spotify.SimplePlaylist
    offset := 0
    limit := 20

    for {
        playlists, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{Offset: &offset, Limit: &limit})
        if err != nil || playlists == nil || len(playlists.Playlists) < 1 {
            break
        }

        for _, playlist := range playlists.Playlists {
            ret = append(ret, playlist)
        }
        // Increment offset for pagination
        offset += limit
    }
    return ret
}

// Returns up to LIMIT tracks from a playlist or album
func FetchTracks(id spotify.ID, limit int, opt string) []spotify.SimpleTrack {
    var ret []spotify.SimpleTrack
    if limit > 0 {
        switch opt {
        case "playlist":
            page, err := client.GetPlaylistTracksOpt(id, &spotify.Options{Limit: &limit}, "")
            if err != nil {
                panic(err)
            }
            for _, track := range page.Tracks {
                ret = append(ret, track.Track.SimpleTrack)
            }
        case "album":
            page, err := client.GetAlbumTracksOpt(id, &spotify.Options{Limit: &limit})
            if err != nil {
                panic(err)
            }
            ret = page.Tracks[:limit]
        }
    }
    return ret
}

// Returns the User object of active client
func GetUser() *spotify.PrivateUser {
    user, err := client.CurrentUser()
    if err != nil {
        panic(err)
    }
    return user
}

// Returns ID of User object of active client
func GetUserID() string {
    return GetUser().ID
}

func GetUserName() string {
    return GetUser().DisplayName
}