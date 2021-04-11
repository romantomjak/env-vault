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

	// set ENVCOMPOSE_PASSWORD_FILE to automatically search for password in that file
	// create to create new encrypted file
	// encrypt to encryp existing file
	// view encrypted file
	// view single secret - for example when executing from scripts, e.g. env-vault -vault secrets.env view -single some.path.to.secret
	// edit encrypted file
	// rekey to change password of encrypted file
	// decrypt to permanently decrypt
}
