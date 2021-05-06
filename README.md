# Spotify Filesystem
The idea is to emulate a simple Unix-like filesystem but for a user's Spotify library. The user should be able to (ideally) view, edit, and create playlists using Linux commands, just like how you would on your Unix-based machines.

## Installation
This project requires `go1.16.3`. To install, simply run:
```
go get github.com/bjma/spotify-filesys
```

## Config
The Spotify API doesn't support playlist folders, so in order to parse folders within your library correctly, create a `.config.json` file, which should look something like this:

```json
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
### `ls`
Displays current folder/playlist and children in specified format; by default, just displays in column/row order (not yet implemented)

Flags:
* `ls --la`
    * Shows playlist status (collaborative/private/public)
* `ls --tree`
    * Shows `pwd` as `tree`
