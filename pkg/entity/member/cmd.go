package member

import "github.com/spf13/cobra"

func ListCmd() *cobra.Command {
	longDoc := `List all the current members in either your organization or your active project.
Either pass the organization flag or the project flag.`
	members := &cobra.Command{
		Use:   "members",
		Short: "List all members of your organization or active project",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			printer.Print(list(cmd))
		},
		Args: cobra.NoArgs,
	}
	flags := members.Flags()
	flags.Bool(organizationFlag, false, "")
	flags.Bool(projectFlag, false, "")
	return members
}
