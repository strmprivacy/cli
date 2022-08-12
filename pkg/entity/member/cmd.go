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

func GetCmd() *cobra.Command {
	member := &cobra.Command{
		Use:   "user",
		Short: "Get a member of your organization",
		PreRun: func(cmd *cobra.Command, args []string) {
			printer = configurePrinter(cmd)
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			printer.Print(get(*&args[0]))
		},
		Args: cobra.ExactArgs(1),
	}
	return member
}

func ManageCmd() *cobra.Command {
	longDoc := `Changes the user roles for the given users. 
All users in the request get all roles specified. Users are specified by their email address.
Possible roles: admin, approver, project-admin, member`
	members := &cobra.Command{
		Use:   "user-roles",
		Short: "Change user roles.",
		Long:  longDoc,
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			manage(cmd)
		},
		Args: cobra.NoArgs,
	}
	flags := members.Flags()
	flags.StringSliceP(usersFlag, "u", []string{}, "Users by email for which roles should be changed")
	flags.StringSliceP(rolesFlag, "r", []string{}, "All roles all users will get")
	return members
}
