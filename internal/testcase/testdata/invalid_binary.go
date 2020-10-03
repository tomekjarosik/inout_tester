package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "this is text on Stderr")
	os.Exit(1)
}
