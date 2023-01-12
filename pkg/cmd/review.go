package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_contract"
)

var ReviewCmd = &cobra.Command{
	Use:               common.ReviewCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Set the state of an entity to IN_REVIEW",
}

func init() {
	ReviewCmd.AddCommand(data_contract.ReviewCmd())
}
