package web_socket

import (
	"github.com/spf13/cobra"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/entity/stream"
)

var longDoc = `The global ` + "`listen`" + ` command is used for starting a Web Socket listener for a stream and output all events to the
console.

This command can receive events from both Source Streams and Derived Streams.

### Usage`

var WebSocketCmd = &cobra.Command{
	Use:               "web-socket (stream-name)",
	Short:             "Read events via the web-socket (not for production purposes)",
	DisableAutoGenTag: true,
	Long:              longDoc,
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, &args[0])
	},
	Args:              cobra.ExactArgs(1), // the stream name
	ValidArgsFunction: stream.NamesCompletion,
}

func init() {
	flags := WebSocketCmd.Flags()
	flags.String(common.ClientIdFlag, "", "Client id to be used for receiving data")
	flags.String(common.ClientSecretFlag, "", "Client secret to be used for receiving data")
}
