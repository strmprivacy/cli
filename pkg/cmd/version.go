package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var GitSha = "dev"
var Version = "dev"
var BuiltOn = "unknown"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		println(fmt.Sprintf("Stream Machine CLI version: %s", Version))
		println(fmt.Sprintf("git commit: %s", GitSha))
		println(fmt.Sprintf("Built on: %s", BuiltOn))
	},
}