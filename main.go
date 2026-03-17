package main

import (
	"fmt"
	"os"

	"github.com/pxblo-x1/devtool-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
