# Spotify Filesystem
The idea is to emulate a simple Unix-like filesystem but for a user's Spotify library. The user should be able to (ideally) view, edit, and create playlists using Linux commands, just like how you would on your Unix-based machines.

## Installation
This project requires `go1.16.3`. To install, simply run:
```
$ go get github.com/bjma/spotify-filesys
```

## Setup
Before doing anything, you should first login to your [dashboard](https://developer.spotify.com/dashboard/login) for Spotify Developers, or create account if you don't have one. There, click the "Create An App" options and set the **Redirect URI** to `http://localhost:8888/callback` in your App settings.

The Spotify API doesn't support playlist folders, so in order to parse folders within your library correctly, you need to create a `.config.json` file, which should look something like this:

```json
{
    "client_id": "your_client_id",
    "client_secret": "your_client_secret",
    "folders": [
        {
            "uri": "spotify:user:user_id:folder:folder_id", 
            "type": "folder", 
            "children": [
                {"type": "playlist", "uri": "spotify:playlist:playlist_id"}, 
                {"type": "playlist", "uri": "spotify:playlist:playlist_id"}, 
                {"type": "playlist", "uri": "spotify:playlist:playlist_id"}, 
            ], 
            "name": "my_folder"
        }
    ]
}
```

You can install the [spotify-folders](https://github.com/mikez/spotify-folders) tool on GitHub to aid with this process; simply drag your playlist into your terminal with the following command typed out:

```
$ spotifyfolders
```

## Supported Commands
### `whoami`
Shows user ID
### `pwd` 
Displays name of current directory (not yet implemented)
### `tree`
Option to show the first `N` tracks (`1 <= N <= 10`)
Support different trees, like
* `--artists`
    * By default, shows Top 10 artists and up to 5 albums
    * Fields:
        * `--artists="artist name"`
* `--albums` (not yet implemented)
    * By default, shows liked albums with artist as parent directory
    * Fields:
        * `--albumns="artist name"`
* `--depth` (kinda buggy still)
    * Prints a tree with number of tracks equal to depth
### `ls`
Displays current folder/playlist and children in specified format; by default, just displays in column/row order (not yet implemented)

Flags:
* `ls --la`
    * Shows playlist status (collaborative/private/public)
* `ls --tree`
    * Shows `pwd` as `tree`