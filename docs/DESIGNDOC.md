A couple concerns:
* How do we deal with authentication?
* Do we need a server?
    * Must have a connection to internet in order to make API requests
* Rate limit: yikes... this might be a problem
* Navigation
    * Need `pwd`
        * Can probably use some field to hold reference to API response? This means the state will always be the most recently requested content
    * Honestly, we can just pull the user's entire playlist library and just create a tree data structure to emulate the filesystem
        * Navigation can then just be done through recursive tree functions
        * Probably need some binary search functions as well
* Persisting OAuth tokens + session?
    * Possibly just store in config as constant and we can just access it anywhere
    * Honestly i think storing the auth tokens in files (like a cache) and reading from them would help everytime we launch
    * If user not authenticated,
        * prompt authentication
    * Else, 
        * read tokens from file and authenticate

## Packages
* [flag](https://pkg.go.dev/flag)

## Source
[tutorial](https://www.rapid7.com/blog/post/2016/08/04/build-a-simple-cli-tool-with-golang/)
[spotify wrapper docs](https://pkg.go.dev/github.com/zmb3/spotify)