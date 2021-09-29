package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/sim"
	"streammachine.io/strm/pkg/sim/randomsim"
)

var SimCmd = &cobra.Command{
	Use:   "sim",
	Short: "Simulate events",
}

func init() {
	flags := SimCmd.PersistentFlags()
	flags.String(sim.SchemaFlag, "streammachine/demo/1.0.2", "what schema to simulate")
	_ = SimCmd.RegisterFlagCompletionFunc(sim.SchemaFlag, schemaCompletion)
	SimCmd.AddCommand(randomsim.RunCmd())
}

func schemaCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	keys := make([]string, 0, len(randomsim.EventGenerators))
	for key := range randomsim.EventGenerators {
		keys = append(keys, key)
	}
	return keys, cobra.ShellCompDirectiveDefault
}
