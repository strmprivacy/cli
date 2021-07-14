package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/sims/randomsim"
)

// SimCmd represents the create command
var SimCmd = &cobra.Command{
	Use:   "sim",
	Short: "Simulate events",
}

func init() {
	SimCmd.AddCommand(randomsim.RunCmd())
}
