package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/diagnostics"
)

var EvaluateCmd = &cobra.Command{
	Use:               common.EvaluateCommandName,
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Calculate Privacy Metrics",
}

func init() {
	EvaluateCmd.AddCommand(diagnostics.EvaluateCmd())
}
