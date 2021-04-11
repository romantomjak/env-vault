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

	// view single secret - for example when executing from scripts, e.g. env-vault -vault secrets.env view -single some.path.to.secret

	// encrypt to encryp existing file
	// decrypt to permanently decrypt
}
