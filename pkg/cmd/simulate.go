package cmd

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/auth"
	"strmprivacy/strm/pkg/simulator/random_events"
)

var SimulateCmd = &cobra.Command{
	Use:               "simulate",
	PersistentPreRun:  auth.RequireAuthenticationPreRun,
	DisableAutoGenTag: true,
	Short:             "Send simulated events with a predefined demo schema (not for production purposes)",
}

func init() {
	flags := SimulateCmd.PersistentFlags()
	flags.String(random_events.SchemaFlag, "strmprivacy/demo/1.0.2", "Which schema to use for creating simulated events")
	_ = SimulateCmd.RegisterFlagCompletionFunc(random_events.SchemaFlag, schemaCompletion)
	SimulateCmd.AddCommand(random_events.RunCmd())
}

func schemaCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	keys := make([]string, 0, len(random_events.EventGenerators))
	for key := range random_events.EventGenerators {
		keys = append(keys, key)
	}
	return keys, cobra.ShellCompDirectiveDefault
}
