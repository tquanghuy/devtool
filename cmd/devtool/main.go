package main

import (
	"flag"
	"fmt"
	"os"

	"devtool/pkg/app"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	a, err := app.Setup(*debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup app: %v\n", err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "app error: %v\n", err)
		os.Exit(1)
	}
}
