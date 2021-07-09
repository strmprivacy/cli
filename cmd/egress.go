package cmd

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/egress"
	"streammachine.io/strm/entity/stream"
	"streammachine.io/strm/sims"
)

// SimCmd represents the create command
var EgressCmd = &cobra.Command{
	Use:               "egress",
	Short:             "Read from egress",
	Run:               func(cmd *cobra.Command, args []string) { egress.Run(cmd, &args[0]) },
	Args:              cobra.ExactArgs(1), // the stream name
	ValidArgsFunction: stream.ExistingNamesCompletion,
}

func init() {
	flags := EgressCmd.Flags()
	flags.String(egress.UrlFlag, "wss://out.dev.strm.services/ws",
		"where to retrieve the events")
	flags.String(sims.ClientIdFlag, "", "client id to be used for sending data")
	flags.String(sims.ClientSecretFlag, "", "client secret to be used for sending data")

}
