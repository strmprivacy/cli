package organization

import (
	"github.com/spf13/cobra"
)

const (
	userEmailsFileFlag = "user-emails-file"
)

var inviteLongDoc = `Invite one or more users to your organization, by email.

Either provide the emails comma-separated on the command line, or pass a file
with the -f flag containing one email address per line.

### Usage`

func InviteUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "users (email,...)",
		Short:             "Invite users to your organization by email",
		Long:              inviteLongDoc,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			inviteUsers(args, cmd)
		},
		Args: cobra.MaximumNArgs(1),
	}

	flags := cmd.Flags()
	flags.StringP(userEmailsFileFlag, "f", "", "file with users to invite, one email per line")
	return cmd
}
