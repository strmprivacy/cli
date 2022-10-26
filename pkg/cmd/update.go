package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/policy"
)

var UpdateCmd = &cobra.Command{
	Use:               common.UpdateCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Update an entity",
}

func init() {
	UpdateCmd.AddCommand(policy.UpdateCmd())
}
