package member

import "github.com/spf13/cobra"

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "project-members [project-name]",
		Short: "List all members of a project. Defaults to active project",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			printer.Print(list(args))
		},
		Args: cobra.MaximumNArgs(1),
	}
}
