package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/sim"
	"streammachine.io/strm/pkg/sim/randomsim"
)

func SimCmd() (cmd *cobra.Command) {
	// SimCmd represents the create command
	simCmd := &cobra.Command{
		Use:   "sim",
		Short: "Simulate events",
	}

	flags := simCmd.PersistentFlags()
	flags.String(sim.SchemaFlag, "streammachine/demo/1.0.2", "what schema to simulate")
	_ = simCmd.RegisterFlagCompletionFunc(sim.SchemaFlag, schemaCompletion)
	simCmd.AddCommand(randomsim.RunCmd())
	return simCmd
}

func schemaCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	keys := make([]string, 0, len(randomsim.EventGenerators))
	for key := range randomsim.EventGenerators {
		keys = append(keys, key)
	}
	return keys, cobra.ShellCompDirectiveDefault
}
