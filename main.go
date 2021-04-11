package main

import (
	"fmt"
	"os"

	"github.com/romantomjak/env-vault/command"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
