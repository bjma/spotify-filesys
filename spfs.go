package main

// Imported packages
import (
	"fmt"
	"flag"

	// Commands
	"github.com/bjma/spotify-filesys/cmd"
)

func main() {
	// Commands
	//textPtr := flag.String("text", "", "Text to parse")
	//metricPtr := flag.String("metric", "chars", "Metric {chars|words}")
	//uniquePtr := flag.Bool("unique", false, "Measure unique values")
	flag.Parse()

	//fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n", *textPtr, *metricPtr, *uniquePtr)
	fmt.Printf(cmd.Hello())
}
