package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/data_contract"
)

var ArchiveCmd = &cobra.Command{
	Use:               common.ArchiveCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Set the state of an entity to ARCHIVED",
}

func init() {
	ArchiveCmd.AddCommand(data_contract.ArchiveCmd())
}
