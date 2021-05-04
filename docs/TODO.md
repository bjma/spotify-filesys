* Figure out a way to make HTTP requests to API from other modules
* Generate a directory tree after establishing TCP connection
* Create a Makefile
* Figure out how to create an interactive shell session
    * Should be able to call `sptfs` to start the shell (like `mininet`)
        * [Possible source](https://hackernoon.com/today-i-learned-making-a-simple-interactive-shell-application-in-golang-aa83adcb266a)
        * [Another](https://sj14.gitlab.io/post/2018/07-01-go-unix-shell/)
        * [This](https://github.com/abiosoft/ishell) package seems really cool
    * `main` should handle authentication and execute the shell
    * Shell should listen for commands (so we have the flags and subcommands here)