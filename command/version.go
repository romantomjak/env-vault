package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/romantomjak/env-vault/version"
)

var versionCmd = &cobra.Command{
	Use:                   "version",
	Short:                 "Prints the env-vault version",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.FullVersion())
	},
}
