package project

import "github.com/spf13/cobra"

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "projects",
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

func CreateCmd() *cobra.Command {
	project := &cobra.Command{
		Use:   "project [name]",
		Short: "Create a new project",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			create(&args[0], cmd)
		},
		Args: cobra.ExactArgs(1), //
	}
	flags := project.Flags()
	flags.String(descriptionFlag, "", "description of the project")
	return project
}

func ManageCmd() *cobra.Command {
	project := &cobra.Command{
		Use:   "project [project-name]",
		Short: "Manage a project, add or remove members. Defaults to active project",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			manage(&args[0], cmd)
		},
		Args: cobra.MaximumNArgs(1),
	}
	flags := project.Flags()
	flags.StringArray(addMembersFlag, []string{}, "[email1,email2,..]")
	flags.StringArray(removeMembersFlag, []string{}, "[email1,email2,..]")
	return project
}
