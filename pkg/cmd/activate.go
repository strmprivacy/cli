package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_contract"
	"strmprivacy/strm/pkg/entity/policy"
)

var ActivateCmd = &cobra.Command{
	Use:               common.ActivateCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Set the state of an entity to ACTIVATED",
}

func init() {
	ActivateCmd.AddCommand(data_contract.ActivateCmd())
	ActivateCmd.AddCommand(policy.ActivateCmd())
}
