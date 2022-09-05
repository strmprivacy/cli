package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/organization"
)

var InviteCmd = &cobra.Command{
	Use:               common.InviteCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Invite users to your organization",
}

func init() {
	InviteCmd.AddCommand(organization.InviteUsersCmd())
}
