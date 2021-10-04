package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/simulator"
	"streammachine.io/strm/pkg/simulator/random_events"
)

var SimulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Send simulated events with a predefined demo schema (not for production purposes)",
}

func init() {
	flags := SimulateCmd.PersistentFlags()
	flags.String(sim.SchemaFlag, "streammachine/demo/1.0.2", "Which schema to use for creating simulated events")
	_ = SimulateCmd.RegisterFlagCompletionFunc(sim.SchemaFlag, schemaCompletion)
	SimulateCmd.AddCommand(random_events.RunCmd())
}

func schemaCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	keys := make([]string, 0, len(random_events.EventGenerators))
	for key := range random_events.EventGenerators {
		keys = append(keys, key)
	}
	return keys, cobra.ShellCompDirectiveDefault
}
