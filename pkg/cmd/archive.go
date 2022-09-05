package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/event_contract"
	"strmprivacy/strm/pkg/entity/schema"
)

var ArchiveCmd = &cobra.Command{
	Use:               common.ArchiveCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Set the state of an entity to ARCHIVED",
}

func init() {
	ArchiveCmd.AddCommand(event_contract.ArchiveCmd())
	ArchiveCmd.AddCommand(schema.ArchiveCmd())
}
