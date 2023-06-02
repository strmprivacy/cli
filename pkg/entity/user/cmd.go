package user

import (
	"github.com/spf13/cobra"
	"strings"
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(`
List all the current users in either your organization or just those that are member of your active project.

Either set the organization flag or the project flag.
`)

func ListCmd() *cobra.Command {
	members := &cobra.Command{
		Use:   "users",
		Short: "List all users of your organization or members of the active project",
		Long:  longDoc,
		Example: util.DedentTrim(`
strm list users --organization
 EMAIL                     FIRST NAME   LAST NAME   USER ROLES

 [...]@strmprivacy.io      bob          rbac        [MEMBER]
 [...]@strmprivacy.io      Demo         STRM        [ADMIN MEMBER]
`),
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
	flags.Bool(organizationFlag, false, "list the users in your organization")
	flags.Bool(projectFlag, false, "list the members of your project")
	return members
}

func GetCmd() *cobra.Command {
	member := &cobra.Command{
		Use:   "user (email)",
		Short: "Get the details of a user in your organization by their email address",
		Example: util.DedentTrim(`
			strm get user [...]@strmprivacy.io -o json
			{
				"email": "[...]@strmprivacy.io",
				"firstName": "Demo",
				"lastName": "STRM",
				"userRoles": [ "ADMIN", "MEMBER" ]
			}
		`),
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
	longDoc := util.DedentTrim(`
			Changes the user roles for the given users. 
			All users in the request get all roles specified. Users are specified by their email address.

			Possible roles: ` + validRolesString)
	members := &cobra.Command{
		Use:   "user-roles",
		Short: "Change user roles.",
		Long:  longDoc,
		Example: util.DedentTrim(`
			strm manage user-roles --roles approver --users user1@example.org,user2@example.org
		`),
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
	flags.StringSliceP(rolesFlag, "r", []string{}, "All roles the selected users will get")
	return members
}

var validRoles = []string{"admin", "approver", "project-admin", "member"}
var validRolesString = strings.Join(validRoles, ", ")
