package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/user"
	"strmprivacy/strm/pkg/entity/project"
)

var ManageCmd = &cobra.Command{
	Use:               common.ManageCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Manage a project or organization",
}

func init() {
	ManageCmd.AddCommand(project.ManageCmd())
	ManageCmd.AddCommand(user.ManageCmd())
}
