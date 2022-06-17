package project

import "github.com/spf13/cobra"

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "projects",
		Short: "List all projects you have access to",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			printer.Print(ListProjects())
		},
	}
}
