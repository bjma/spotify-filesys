// This is just for testing locality
package cmd

import (
    "fmt"
    "flag"
)

// Have command flag as global so we can reference
var helloCommand *flag.FlagSet

// Each command should take in an array of strings (flags + other args)
func HelloInit(arr []string) {
    helloCommand = flag.NewFlagSet("hello", flag.ExitOnError)
    // Flags
    reverseHelloPtr := helloCommand.Bool("reverse", false, "Reverses greeting message.")
    // Parse flags
    helloCommand.Parse(arr)
    // Matching flags
    if helloCommand.Parsed() {
        if *reverseHelloPtr {
            HelloExecute(true)
        } else {
            HelloExecute(false)
        }
    }
}

// Actual execution of command
func HelloExecute(reverse bool) {
    sentence := [3]string{"Hello", " ", "World"}

    if reverse {
        for i := len(sentence) - 1; i >= 0; i-- {
            fmt.Print(sentence[i])
        }
    } else {
        for i := 0; i < len(sentence); i++ {
            fmt.Print(sentence[i])
        }
    }
    fmt.Printf("\n")
}