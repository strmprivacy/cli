package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/organization"
)

var InviteCmd = &cobra.Command{
	Use:               common.InviteCommandName,
	DisableAutoGenTag: true,
	Short:             "Invite users to your organization",
}

func init() {
	InviteCmd.AddCommand(organization.InviteUsersCmd())
}
