package web_socket

import (
	"github.com/spf13/cobra"
	"streammachine.io/strm/pkg/entity/stream"
	"streammachine.io/strm/pkg/simulator"
)

var WebSocketCmd = &cobra.Command{
	Use:   "web-socket (stream-name)",
	Short: "Read events via the web-socket (not for production purposes)",
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, &args[0])
	},
	Args:              cobra.ExactArgs(1), // the stream name
	ValidArgsFunction: stream.NamesCompletion,
}

func init() {
	flags := WebSocketCmd.Flags()
	flags.String(sim.ClientIdFlag, "", "Client id to be used for receiving data")
	flags.String(sim.ClientSecretFlag, "", "Client secret to be used for receiving data")
}
