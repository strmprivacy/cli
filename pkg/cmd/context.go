package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var ContextCommand = &cobra.Command{
	Use:   "context",
	Short: "Print the CLI context",
	Run: func(cmd *cobra.Command, _ []string) {
		printContext(cmd)
	},
}

func printContext(cmd *cobra.Command) {
	fmt.Println(cmd.Context())
}
