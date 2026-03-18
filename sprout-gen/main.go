package main

import (
	"os"

	"github.com/ink-yht-code/sprout/sprout-gen/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
