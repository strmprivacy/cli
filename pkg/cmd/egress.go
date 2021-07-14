package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/egress"
	"streammachine.io/strm/pkg/entity/stream"
	"streammachine.io/strm/pkg/sims"
)

// SimCmd represents the create command
var EgressCmd = &cobra.Command{
	Use:               "egress",
	Short:             "Read from egress",
	Run:               func(cmd *cobra.Command, args []string) {
		egress.Run(cmd, &args[0])
	},
	Args:              cobra.ExactArgs(1), // the stream name
	ValidArgsFunction: stream.NamesCompletion,
}

func init() {
	flags := EgressCmd.Flags()
	flags.String(sims.ClientIdFlag, "", "Client id to be used for receiving data")
	flags.String(sims.ClientSecretFlag, "", "Client secret to be used for receiving data")

}
