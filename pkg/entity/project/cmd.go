package project

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/projects/v1"
	"strmprivacy/strm/pkg/common"
)

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "projects",
		Short: "List all projects you have access to",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			printer.Print(ListProjectsWithActive())
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
		Short: "Manage a project: add or remove members. Defaults to active project",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			manage(args, cmd)
		},
		Args: cobra.MaximumNArgs(1),
	}
	flags := project.Flags()
	flags.StringArray(addMembersFlag, []string{}, "[email1,email2,..]")
	flags.StringArray(removeMembersFlag, []string{}, "[email1,email2,..]")
	return project
}

func GetCmd() *cobra.Command {
	project := &cobra.Command{
		Use:   "project [project-name]",
		Short: "Get a project",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			printer.Print(get(args[0]))
		},
		Args: cobra.ExactArgs(1),
	}
	return project
}

func DeleteCmd() *cobra.Command {
	project := &cobra.Command{
		Use:   "project [project-name]",
		Short: "Delete a project and all associated resources",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			del(args[0])
		},
		Args: cobra.ExactArgs(1),
	}
	flags := project.Flags()
	flags.StringArray(addMembersFlag, []string{}, "[email1,email2,..]")
	flags.StringArray(removeMembersFlag, []string{}, "[email1,email2,..]")
	return project
}

func NamesCompletion(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	req := &projects.ListProjectsRequest{}
	response, err := client.ListProjects(apiContext, req)

	if err != nil {
		return common.GrpcRequestCompletionError(err)
	}

	names := make([]string, 0, len(response.Projects))
	for _, p := range response.Projects {
		names = append(names, p.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
