package cmd

import (
	"github.com/spf13/cobra"
	"github.com/strmprivacy/api-definitions-go/v2/api/monitoring/v1"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/monitor"
)

var MonitorCmd = &cobra.Command{
	Use:               common.MonitorCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Monitor STRM entities",
}

func init() {
	MonitorCmd.AddCommand(monitor.Command(monitoring.EntityState_BATCH_EXPORTER))
	MonitorCmd.AddCommand(monitor.Command(monitoring.EntityState_BATCH_JOB))
	MonitorCmd.AddCommand(monitor.Command(monitoring.EntityState_STREAM))
}
