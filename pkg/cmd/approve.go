package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_contract"
)

var ApproveCmd = &cobra.Command{
	Use:               common.ApproveCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Set the state of an entity to APPROVED",
}

func init() {
	ApproveCmd.AddCommand(data_contract.ApproveCmd())
}
