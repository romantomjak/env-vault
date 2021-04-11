package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// GitCommit that was compiled. This will be filled in by the compiler
	GitCommit string

	// Version number that is being run at the moment.
	// This will be filled in by the compiler
	Version string
)

var versionCmd = &cobra.Command{
	Use:                   "version",
	Short:                 "Prints the env-vault version",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("env-vault v%s (%s)\n", Version, GitCommit)
	},
}
