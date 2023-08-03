package cmd

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v3/api/monitoring/v1"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/logs"
)

var LogsCmd = &cobra.Command{
	Use:               common.LogsCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Show logs of STRM entities",
}

func init() {
	LogsCmd.AddCommand(logs.Command(monitoring.EntityState_BATCH_EXPORTER))
	LogsCmd.AddCommand(logs.Command(monitoring.EntityState_BATCH_JOB))
	LogsCmd.AddCommand(logs.Command(monitoring.EntityState_STREAM))
}
